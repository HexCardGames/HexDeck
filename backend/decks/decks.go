package decks

import (
	"github.com/HexCardGames/HexDeck/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func DeckFromInterface(cardDeckId int, cardDeck bson.D) types.CardDeck {
	bsonBytes, _ := bson.Marshal(cardDeck)

	switch cardDeckId {
	case 0:
		deck := Classic{}
		bson.Unmarshal(bsonBytes, &deck)
		return &deck
	case 1:
		deck := HexV1{}
		bson.Unmarshal(bsonBytes, &deck)
		return &deck
	}

	return nil
}

func CardFromInterface(cardDeckId int, card bson.D) types.Card {
	bsonBytes, _ := bson.Marshal(card)

	switch cardDeckId {
	case 0:
		deck := ClassicCard{}
		bson.Unmarshal(bsonBytes, &deck)
		return &deck
	case 1:
		deck := HexV1Card{}
		bson.Unmarshal(bsonBytes, &deck)
		return &deck
	}
	return nil
}
