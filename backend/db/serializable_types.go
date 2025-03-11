package db

import (
	"sync"

	"github.com/HexCardGames/HexDeck/decks"
	"github.com/HexCardGames/HexDeck/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SerializablePlayer struct {
	PlayerId     bson.ObjectID
	SessionToken string
	Username     string
	Permissions  int
	Cards        []bson.D
}

func (serializable *SerializablePlayer) ToPlayer(cardDeckId int) types.Player {
	cards := make([]types.Card, len(serializable.Cards))
	for i, card := range serializable.Cards {
		cards[i] = decks.CardFromInterface(cardDeckId, card)
	}
	player := types.Player{
		PlayerId:     serializable.PlayerId,
		SessionToken: serializable.SessionToken,
		Username:     serializable.Username,
		Permissions:  serializable.Permissions,
		Connection:   types.WebsocketConnection{IsConnected: false},
		Cards:        cards,
		Mutex:        &sync.Mutex{},
	}
	player.ResetInactivity()
	return player
}

type SerializableRoom struct {
	RoomId      bson.ObjectID `bson:"_id"`
	JoinCode    string
	GameState   types.GameState
	GameOptions types.GameOptions
	CardDeckId  int
	CardDeck    bson.D
	Players     []SerializablePlayer
	OwnerId     bson.ObjectID
	MoveTimeout int
	Winner      *bson.ObjectID
}

func (serializable SerializableRoom) ToRoom() *types.Room {
	players := make([]*types.Player, len(serializable.Players))
	for i, serializablePlayer := range serializable.Players {
		player := serializablePlayer.ToPlayer(serializable.CardDeckId)
		players[i] = &player
	}
	cardDeck := decks.DeckFromInterface(serializable.CardDeckId, serializable.CardDeck)
	room := &types.Room{
		RoomId:       serializable.RoomId,
		JoinCode:     serializable.JoinCode,
		GameState:    serializable.GameState,
		GameOptions:  serializable.GameOptions,
		CardDeckId:   serializable.CardDeckId,
		CardDeck:     cardDeck,
		Players:      players,
		PlayersMutex: &sync.Mutex{},
		OwnerId:      serializable.OwnerId,
		MoveTimeout:  serializable.MoveTimeout,
		Winner:       serializable.Winner,
	}
	cardDeck.SetRoom(room)
	return room
}
