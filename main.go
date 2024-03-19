package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mario-so/mteag-cli/api"
	"github.com/mario-so/mteag-cli/item"
)

type Mode int64

const (
	Search int64 = iota
	Select
	View
)

type model struct {
	textInput    textinput.Model
	cardList     list.Model
	infoTable    table.Model
	spinner      spinner.Model
	selectedCard item.Card
	mode         Mode
	isLoading    bool
}

type getCardMsg struct {
	cards []list.Item
}

func getCardsCmd(cardName string) tea.Cmd {
	return func() tea.Msg {
		cards, err := api.GetCards(cardName)
		if err != nil {
			return nil
		}
		cardListItems := make([]list.Item, len(cards))
		for i, card := range cards {
			cardListItems[i] = &ui.CardListItem{Card: card}
		}
		return getCardsMsg{cards: cardListItems}
	}
}

func main() {
	// Get details for a specific card
	card, _ := GetCards("jace sculptor")

	// Print the details
	log.Println(card.Name, card.ManaCost, card.ImageUris.Large, card.Prices.Usd, card.Prices.Eur)
}
