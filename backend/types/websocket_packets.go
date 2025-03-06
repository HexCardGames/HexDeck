package types

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type S2C_Status struct {
	IsError    bool
	StatusCode string
	Message    string
}
type S2C_PlayerInfo struct {
	PlayerId    bson.ObjectID
	Username    string
	Permissions int
	IsConnected bool
}
type S2C_RoomInfo struct {
	RoomId      bson.ObjectID `bson:"_id"`
	JoinCode    string
	GameState   GameState
	GameOptions GameOptions
	TopCard     Card
	CardDeckId  int
	Winner      *bson.ObjectID
	Players     []S2C_PlayerInfo
}
type S2C_Card struct {
	CanPlay bool
	Card    Card
}
type S2C_OwnCards struct {
	Cards []S2C_Card
}
type S2C_PlayerState struct {
	PlayerId bson.ObjectID
	NumCards int
	Active   bool
}
type S2C_CardPlayed struct {
	Card      Card
	CardIndex int
	PlayedBy  bson.ObjectID
}
type S2C_PlayedCardUpdate struct {
	UpdatedBy bson.ObjectID
	Card      Card
}

type C2S_UpdatePlayer struct {
	PlayerId    bson.ObjectID
	Username    *string
	Permissions *int
}
type C2S_KickPlayer struct {
	PlayerId bson.ObjectID
}
type C2S_PlayCard struct {
	CardIndex *int
	CardData  interface{}
}
type C2S_UpdatePlayedCard struct {
	CardData interface{}
}

func BuildRoomInfoPacket(room *Room) S2C_RoomInfo {
	players := make([]S2C_PlayerInfo, len(room.Players))
	for i, player := range room.Players {
		players[i] = S2C_PlayerInfo{
			PlayerId:    player.PlayerId,
			Username:    player.Username,
			Permissions: player.Permissions,
			IsConnected: player.Connection.IsConnected,
		}
	}
	roomInfo := S2C_RoomInfo{
		RoomId:      room.RoomId,
		JoinCode:    room.JoinCode,
		GameState:   room.GameState,
		CardDeckId:  room.CardDeckId,
		GameOptions: room.GameOptions,
		Winner:      room.Winner,
		Players:     players,
	}

	if room.CardDeck != nil {
		roomInfo.TopCard = room.CardDeck.GetTopCard()
	}
	return roomInfo
}

func BuildOwnCardsPacket(room *Room, player *Player) S2C_OwnCards {
	cards := make([]S2C_Card, len(player.Cards))
	for i, card := range player.Cards {
		cards[i] = S2C_Card{
			Card:    card,
			CanPlay: room.CardDeck.CanPlay(card),
		}
	}
	return S2C_OwnCards{
		Cards: cards,
	}
}

func BuildPlayerStatePacket(room *Room, player *Player) S2C_PlayerState {
	return S2C_PlayerState{PlayerId: player.PlayerId, NumCards: len(player.Cards), Active: room.CardDeck.IsPlayerActive(player)}
}
func BuildCardPlayedPacket(player *Player, cardIndex int, card Card) S2C_CardPlayed {
	return S2C_CardPlayed{Card: card, CardIndex: cardIndex, PlayedBy: player.PlayerId}
}
func BuildPlayedCardUpdatePacket(player *Player, card Card) S2C_PlayedCardUpdate {
	return S2C_PlayedCardUpdate{UpdatedBy: player.PlayerId, Card: card}
}
