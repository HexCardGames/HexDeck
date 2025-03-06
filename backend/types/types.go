package types

import (
	"sync"

	"github.com/zishang520/socket.io/v2/socket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WebsocketConnection struct {
	IsConnected bool
	Socket      *socket.Socket
}

type Card interface {
}

type CardDeck interface {
	Init(*Room)
	SetRoom(*Room)
	IsEmpty() bool
	DrawCard() Card
	CanPlay(Card) bool
	PlayCard(Card) bool
	GetTopCard() Card
	UpdatePlayedCard(interface{}) Card
	IsPlayerActive(*Player) bool
}

type Player struct {
	PlayerId          bson.ObjectID
	SessionToken      string
	Username          string
	Permissions       int
	Cards             []Card              `json:"-"`
	Connection        WebsocketConnection `bson:"-" json:"-"`
	InactivityTimeout int                 `bson:"-" json:"-"`
}

func (player *Player) ResetInactivity() {
	player.InactivityTimeout = 20 * 1000
}

func (player *Player) SetPermissionBit(bit RoomPermission) {
	player.Permissions |= (1 << bit)
}

func (player *Player) ClearPermissionBit(bit RoomPermission) {
	player.Permissions &= ^(1 << bit)
}

func (player *Player) HasPermissionBit(bit RoomPermission) bool {
	return player.Permissions&(1<<bit) > 0
}

type GameState int

const (
	StateLobby GameState = iota
	StateRunning
	StateEnded
)

type RoomPermission int

const (
	PermissionHost RoomPermission = 0
)

type GameOptions struct {
}

type Room struct {
	RoomId       bson.ObjectID `bson:"_id"`
	JoinCode     string
	GameState    GameState
	GameOptions  GameOptions
	CardDeckId   int
	CardDeck     CardDeck
	Players      []*Player
	PlayersMutex *sync.Mutex `bson:"-"`
	OwnerId      bson.ObjectID
	MoveTimeout  int
	Winner       *bson.ObjectID
}

func (room *Room) AppendPlayer(player *Player) {
	room.PlayersMutex.Lock()
	defer room.PlayersMutex.Unlock()
	room.Players = append(room.Players, player)
}

func (room *Room) RemovePlayer(target Player) bool {
	room.PlayersMutex.Lock()
	defer room.PlayersMutex.Unlock()
	return room.RemovePlayerUnsafe(target)
}

func (room *Room) FindPlayer(playerId bson.ObjectID) *Player {
	room.PlayersMutex.Lock()
	defer room.PlayersMutex.Unlock()

	for _, player := range room.Players {
		if player.PlayerId == playerId {
			return player
		}
	}
	return nil
}

func (room *Room) RemovePlayerUnsafe(target Player) bool {
	foundHost := false
	foundPlayer := false
	for i := 0; i < len(room.Players); i++ {
		player := room.Players[i]
		if player.PlayerId == target.PlayerId {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			foundPlayer = true
			i--
			continue
		}
		if player.HasPermissionBit(PermissionHost) {
			foundHost = true
		}
	}
	if !foundPlayer {
		return false
	}
	if !foundHost && len(room.Players) > 0 {
		room.Players[0].SetPermissionBit(PermissionHost)
	}
	return true
}

func (room *Room) IsUsernameAvailable(username string) bool {
	for _, player := range room.Players {
		if player.Username == username {
			return false
		}
	}
	return true
}
