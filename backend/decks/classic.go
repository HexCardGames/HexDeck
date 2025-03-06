package decks

import (
	"strconv"

	"github.com/HexCardGames/HexDeck/types"
	"github.com/HexCardGames/HexDeck/utils"
)

type Classic struct {
	room              *types.Room
	CardsPlayed       []*ClassicCard
	CardsRemaining    []*ClassicCard
	DirectionReversed bool
	ActivePlayer      int
}

var ClassicColors = []string{"red", "yellow", "blue", "green"}

func (deck *Classic) Init(room *types.Room) {
	deck.room = room
	deck.DirectionReversed = false
	deck.ActivePlayer = 0
	cards := make([]*ClassicCard, 108)
	offset := 0
	for i := 0; i < 4; i++ {
		color := ClassicColors[i]
		for j := 0; j < 19; j++ {
			cards[offset] = &ClassicCard{
				Symbol: strconv.Itoa((j + 1) % 10),
				Color:  color,
			}
			offset += 1
		}
		for j := 0; j < 2; j++ {
			cards[offset] = &ClassicCard{Symbol: "action:skip", Color: color}
			cards[offset+1] = &ClassicCard{Symbol: "action:reverse", Color: color}
			cards[offset+2] = &ClassicCard{Symbol: "action:draw_2", Color: color}
			offset += 3
		}
		cards[offset] = &ClassicCard{Symbol: "action:wildcard", Color: "black"}
		cards[offset+1] = &ClassicCard{Symbol: "action:draw_4", Color: "black"}
		offset += 2
	}
	utils.ShuffleSlice(&cards)
	deck.CardsRemaining = cards
}

func (deck *Classic) SetRoom(room *types.Room) {
	deck.room = room
}

func (deck *Classic) IsEmpty() bool {
	return len(deck.CardsRemaining) == 0
}

func (deck *Classic) getTopCard() *ClassicCard {
	if len(deck.CardsPlayed) == 0 {
		return nil
	}
	return deck.CardsPlayed[len(deck.CardsPlayed)-1]
}

func (deck *Classic) GetTopCard() types.Card {
	return deck.getTopCard()
}

func (deck *Classic) drawCard(player *types.Player) types.Card {
	if deck.IsEmpty() {
		return nil
	}

	card := deck.CardsRemaining[0]
	deck.CardsRemaining = deck.CardsRemaining[1:]
	player.Cards = append(player.Cards, card)
	return card
}

func (deck *Classic) getActivePlayer() int {
	return utils.Mod(deck.ActivePlayer, len(deck.room.Players))
}

func (deck *Classic) DrawCard() types.Card {
	// Can't draw another card before wildcard color is selected
	topCard := deck.getTopCard()
	if topCard != nil && topCard.Color == "black" {
		return nil
	}

	card := deck.drawCard(deck.room.Players[deck.getActivePlayer()])
	deck.nextPlayer()
	return card
}

func (deck *Classic) getNextPlayer() int {
	direction := 1
	if deck.DirectionReversed {
		direction = -1
	}
	return utils.Mod((deck.ActivePlayer + direction), len(deck.room.Players))
}

func (deck *Classic) nextPlayer() {
	deck.ActivePlayer = deck.getNextPlayer()
}

func (deck *Classic) CanPlay(card types.Card) bool {
	topCard := deck.getTopCard()
	checkCard := card.(*ClassicCard)
	if topCard == nil || checkCard == nil {
		return topCard == nil
	}
	return checkCard.Color == "black" || checkCard.Color == topCard.Color || checkCard.Symbol == topCard.Symbol
}

func (deck *Classic) PlayCard(card types.Card) bool {
	if !deck.CanPlay(card) {
		return false
	}
	deckCard := card.(*ClassicCard)
	deck.CardsPlayed = append(deck.CardsPlayed, deckCard)

	if deckCard.Symbol == "action:skip" {
		deck.nextPlayer()
	} else if deckCard.Symbol == "action:draw_2" || deckCard.Symbol == "action:draw_4" {
		targetPlayer := deck.room.Players[deck.getNextPlayer()]
		amount := 2
		if deckCard.Symbol == "action:draw_4" {
			amount = 4
		}
		for range amount {
			card := deck.drawCard(targetPlayer)
			if card == nil {
				// TODO: Handle empty card deck
				break
			}
		}
	} else if deckCard.Symbol == "action:reverse" {
		deck.DirectionReversed = !deck.DirectionReversed
	}

	if deckCard.Color != "black" {
		deck.nextPlayer()
	}

	return true
}

func (deck *Classic) UpdatePlayedCard(cardData interface{}) types.Card {
	topCard := deck.getTopCard()
	if topCard.Color != "black" {
		return nil
	}
	updateData, ok := cardData.(map[string]interface{})
	if !ok {
		return nil
	}
	newColor, ok := updateData["Color"].(string)
	if !ok {
		return nil
	}

	for _, color := range ClassicColors {
		if newColor == color {
			deck.nextPlayer()
			topCard.Color = color
			return topCard
		}
	}
	return nil
}

func (deck *Classic) IsPlayerActive(target *types.Player) bool {
	return deck.room.Players[utils.Mod(deck.ActivePlayer, len(deck.room.Players))] == target
}

type ClassicCard struct {
	Symbol string
	Color  string
}
