package decks

import (
	"fmt"
	"math/rand/v2"

	"github.com/HexCardGames/HexDeck/types"
	"github.com/HexCardGames/HexDeck/utils"
)

type HexV1 struct {
	room        *types.Room
	CardsPlayed []*HexV1Card
	PlayerOrder []int
	ActiveIndex int
}

type HexV1Card struct {
	Symbol       string
	Color        string
	NumericValue int
}

var HexV1Colors = []string{"blue", "green", "yellow", "purple"}
var HexV1ActionCards = []string{"shuffle", "skip", "draw", "swap"}

func (deck *HexV1) Init(room *types.Room) {
	deck.room = room
	deck.PlayerOrder = make([]int, len(room.Players))
	deck.ActiveIndex = 0

	deck.room.PlayersMutex.Lock()
	defer deck.room.PlayersMutex.Unlock()
	for i, player := range deck.room.Players {
		deck.PlayerOrder[i] = i
		player.Mutex.Lock()
		defer player.Mutex.Unlock()
		deck.drawMany(player, 8)
	}
}

func (deck *HexV1) SetRoom(room *types.Room) {
	deck.room = room
}

func (deck *HexV1) IsEmpty() bool {
	return false
}

func (deck *HexV1) getTopCard() *HexV1Card {
	if len(deck.CardsPlayed) == 0 {
		return nil
	}
	return deck.CardsPlayed[len(deck.CardsPlayed)-1]
}

func (deck *HexV1) GetTopCard() types.Card {
	return deck.getTopCard()
}

func (deck *HexV1) generateCard() *HexV1Card {
	cardType := rand.IntN(16 + len(HexV1ActionCards))
	cardColor := HexV1Colors[rand.IntN(len(HexV1Colors))]
	if cardType < 16 {
		return &HexV1Card{
			Symbol:       fmt.Sprintf("%x", cardType),
			Color:        cardColor,
			NumericValue: cardType,
		}
	}
	cardSymbol := HexV1ActionCards[cardType-16]
	if rand.IntN(100) <= 10 {
		cardColor = "rainbow"
	}
	return &HexV1Card{
		Symbol:       "action:" + cardSymbol,
		Color:        cardColor,
		NumericValue: 3,
	}
}

func (deck *HexV1) drawCard(player *types.Player) types.Card {
	card := deck.generateCard()
	player.Cards = append(player.Cards, card)
	return card
}

func (deck *HexV1) drawMany(player *types.Player, cards int) {
	for i := 0; i < cards; i++ {
		deck.drawCard(player)
	}
}

func (deck *HexV1) DrawCard() types.Card {
	// Can't draw another card before wildcard color is selected
	topCard := deck.getTopCard()
	if topCard != nil && topCard.Color == "rainbow" {
		return nil
	}

	card := deck.drawCard(deck.getPlayer(deck.ActiveIndex))
	deck.nextPlayer()
	return card
}

func (deck *HexV1) getNextValidIndex(index int) int {
	if len(deck.room.Players) == 0 || len(deck.PlayerOrder) == 0 {
		return -1
	}
	checkIndex := utils.Mod(index, len(deck.PlayerOrder))
	for deck.PlayerOrder[checkIndex] >= len(deck.room.Players) {
		checkIndex = utils.Mod(checkIndex+1, len(deck.PlayerOrder))
	}
	return checkIndex
}

func (deck *HexV1) getPlayer(index int) *types.Player {
	playerIndex := deck.getNextValidIndex(index)
	if playerIndex == -1 {
		return nil
	}
	return deck.room.Players[deck.PlayerOrder[playerIndex]]
}

func (deck *HexV1) getNextPlayerIndex() int {
	return deck.getNextValidIndex(deck.ActiveIndex + 1)
}

func (deck *HexV1) nextPlayer() {
	deck.ActiveIndex = deck.getNextPlayerIndex()
}

func (deck *HexV1) IsPlayerActive(target *types.Player) bool {
	return deck.getPlayer(deck.ActiveIndex) == target
}

func (deck *HexV1) CanPlay(card types.Card) bool {
	topCard := deck.getTopCard()
	checkCard := card.(*HexV1Card)
	if topCard == nil || checkCard == nil {
		return topCard == nil
	}
	return topCard.Color != "rainbow" && (checkCard.Color == "rainbow" || checkCard.Color == topCard.Color || checkCard.Symbol == topCard.Symbol)
}

func (deck *HexV1) PlayCard(card types.Card) bool {
	if !deck.CanPlay(card) {
		return false
	}
	deckCard := card.(*HexV1Card)
	targetPlayer := deck.getPlayer(deck.ActiveIndex)
	nextPlayer := deck.getPlayer(deck.getNextPlayerIndex())

	if deckCard.Symbol == "action:skip" && deckCard.Color != "rainbow" {
		deck.nextPlayer()
	} else if deckCard.Symbol == "action:draw" {
		amount := 3
		topCard := deck.getTopCard()
		if topCard != nil {
			amount = topCard.NumericValue
		}
		deck.drawMany(nextPlayer, amount)
	} else if deckCard.Symbol == "action:shuffle" {
		utils.ShuffleSlice(&deck.PlayerOrder)
	} else if deckCard.Symbol == "action:swap" {
		p1Cards := targetPlayer.Cards
		p2Cards := nextPlayer.Cards
		targetPlayer.Cards = p2Cards
		nextPlayer.Cards = p1Cards
	}

	if deckCard.Color != "rainbow" {
		deck.nextPlayer()
	}

	deck.CardsPlayed = append(deck.CardsPlayed, deckCard)
	return true
}

func (deck *HexV1) UpdatePlayedCard(cardData interface{}) types.Card {
	topCard := deck.getTopCard()
	if topCard.Color != "rainbow" {
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

	for _, color := range HexV1Colors {
		if newColor == color {
			deck.nextPlayer()
			if topCard.Symbol == "action:skip" {
				deck.nextPlayer()
			}
			topCard.Color = color
			return topCard
		}
	}
	return nil
}
