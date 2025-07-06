package game

import (
	"log/slog"
	"math/rand/v2"
	"strconv"
	"sync"

	"github.com/HexCardGames/HexDeck/db"
	"github.com/HexCardGames/HexDeck/decks"
	"github.com/HexCardGames/HexDeck/types"
	"github.com/HexCardGames/HexDeck/utils"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var roomsMutex sync.Mutex = sync.Mutex{}
var rooms []*types.Room = make([]*types.Room, 0)

func GenerateJoinCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.IntN(10))
	}
	return code
}

func LoadRooms() {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	rooms = db.Conn.QueryRunningRooms()
}

func CreateRoom() *types.Room {
	newRoom := &types.Room{
		RoomId:       bson.NewObjectID(),
		JoinCode:     GenerateJoinCode(),
		GameState:    types.StateLobby,
		Players:      make([]*types.Player, 0),
		PlayersMutex: &sync.Mutex{},
		CardDeckId:   1,
	}

	db.Conn.InsertRoom(newRoom)
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	rooms = append(rooms, newRoom)
	return newRoom
}

func FindRoomByJoinCode(joinCode string) *types.Room {
	for _, room := range rooms {
		if room.JoinCode != joinCode {
			continue
		}
		return room
	}
	return nil
}

func FindSession(sessionToken string) (*types.Room, *types.Player) {
	for _, room := range rooms {
		for _, player := range room.Players {
			if player.SessionToken == sessionToken {
				return room, player
			}
		}
	}
	return nil, nil
}

func JoinRoom(room *types.Room, requestedUsername string) *types.Player {
	var username string
	if requestedUsername != "" && room.IsUsernameAvailable(requestedUsername) {
		username = requestedUsername
	} else {
		username = petname.Generate(2, " ")
	}

	player := &types.Player{
		PlayerId:     bson.NewObjectID(),
		SessionToken: uuid.New().String(),
		Username:     username,
		Permissions:  0,
		Cards:        make([]types.Card, 0),
		Connection: types.WebsocketConnection{
			IsConnected: false,
		},
		Mutex: &sync.Mutex{},
	}
	player.ResetInactivity()
	room.AppendPlayer(player)
	OnRoomUpdate(room)
	return player
}

type GameStats struct {
	RunningGames      int
	OnlinePlayerCount int
}

func CalculateStats() GameStats {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	stats := GameStats{RunningGames: 0, OnlinePlayerCount: 0}
	for _, game := range rooms {
		stats.RunningGames += 1
		for _, player := range game.Players {
			if player.Connection.IsConnected {
				stats.OnlinePlayerCount += 1
			}
		}
	}
	return stats
}

func UpdateGameState(room *types.Room, newState types.GameState) {
	if room.GameState != types.StateEnded && newState == types.StateEnded {
		db.Conn.IncrementGamesPlayed()
	}
	room.GameState = newState
	OnRoomUpdate(room)
}

func SetCardDeck(room *types.Room, id int) bool {
	if id < 0 || id > 1 {
		return false
	}
	room.CardDeckId = id
	OnRoomUpdate(room)
	return true
}

func CreateCardDeckObj(room *types.Room) {
	switch room.CardDeckId {
	case 0:
		room.CardDeck = &decks.Classic{}
	case 1:
		room.CardDeck = &decks.HexV1{}
	}
}

func BroadcastInRoom(room *types.Room, topic string, data interface{}) {
	for _, player := range room.Players {
		if !player.Connection.IsConnected || player.Connection.Socket == nil {
			continue
		}
		player.Connection.Socket.Emit(topic, data)
	}
}

func SendInitialData(room *types.Room, targetPlayer *types.Player) {
	if targetPlayer.Connection.Socket == nil {
		return
	}
	targetPlayer.Connection.Socket.Emit("OwnCards", types.BuildOwnCardsPacket(room, targetPlayer))
	for _, player := range room.Players {
		targetPlayer.Connection.Socket.Emit("PlayerState", types.BuildPlayerStatePacket(room, player))
	}
}

func OnRoomUpdate(room *types.Room) {
	db.Conn.UpdateRoom(room)
	BroadcastInRoom(room, "RoomInfo", types.BuildRoomInfoPacket(room))
}

func OnPlayerStateUpdate(room *types.Room, player *types.Player, skipDBUpdate bool) {
	if !skipDBUpdate {
		db.Conn.UpdateRoom(room)
	}
	if player.Connection.Socket == nil {
		return
	}
	player.Connection.Socket.Emit("OwnCards", types.BuildOwnCardsPacket(room, player))
	BroadcastInRoom(room, "PlayerState", types.BuildPlayerStatePacket(room, player))
}

func UpdateAllPlayers(room *types.Room) {
	db.Conn.UpdateRoom(room)
	for _, player := range room.Players {
		OnPlayerStateUpdate(room, player, true)
	}
}

func OnPlayCard(room *types.Room, player *types.Player, cardIndex int, card types.Card) {
	BroadcastInRoom(room, "CardPlayed", types.BuildCardPlayedPacket(player, cardIndex, card))
	UpdateAllPlayers(room)
}

func OnPlayedCardUpdate(room *types.Room, player *types.Player, card types.Card) {
	BroadcastInRoom(room, "PlayedCardUpdate", types.BuildPlayedCardUpdatePacket(player, card))
	UpdateAllPlayers(room)
}

func StartGame(room *types.Room) {
	if room.GameState != types.StateLobby {
		return
	}
	CreateCardDeckObj(room)
	room.CardDeck.Init(room)
	UpdateGameState(room, types.StateRunning)
	UpdateAllPlayers(room)
}

func TickRooms(deltaTime int) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	for i := 0; i < len(rooms); i++ {
		room := rooms[i]

		hasChanged := false
		room.PlayersMutex.Lock()
		for j := 0; j < len(room.Players); j++ {
			player := room.Players[j]
			if player.Connection.IsConnected {
				continue
			}
			if player.InactivityTimeout <= deltaTime {
				slog.Debug("Removing player from room due to inactivity", "username", player.Username, "playerId", player.PlayerId.Hex(), "roomId", room.RoomId.Hex())
				hasChanged = true
				room.RemovePlayerUnsafe(*player)
				j--
			}
			player.InactivityTimeout -= deltaTime
		}

		if len(room.Players) == 0 {
			slog.Debug("Ending and unloading empty room", "roomId", room.RoomId.Hex())
			UpdateGameState(room, types.StateEnded)
			utils.RemoveSliceElement(&rooms, room)
			i--
		}
		room.PlayersMutex.Unlock()

		if hasChanged {
			OnRoomUpdate(room)
		}
	}
}
