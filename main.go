package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mario-so/mteag-cli/api"
	"github.com/mario-so/mteag-cli/item"
	"github.com/mario-so/mteag-cli/ui"
)

type Mode int64

const (
	Search Mode = iota
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

type getCardsMsg struct {
	cards []list.Item
}

func getCardsCmd(cardName string) tea.Cmd {
	return func() tea.Msg {
		card, err := api.GetCard(cardName) // Fetch a single card
		if err != nil {
			fmt.Printf("Failed to retrieve card details: %v\n", err) // Print with proper error handling
			return nil
		}
		cardListItem := &ui.CardListItem{Card: card}
		cardListItems := []list.Item{cardListItem}

		return getCardsMsg{cards: cardListItems}
	}
}

func (m model) setInfoTable() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Mana Cost", Width: 20},
		{Title: "Price (USD)", Width: 20},
		{Title: "Price (EUR)", Width: 20},
		{Title: "Image", Width: 100},
	}

	rows := []table.Row{
		{m.selectedCard.Name, m.selectedCard.ManaCost, m.selectedCard.Prices.Usd, m.selectedCard.Prices.Eur, m.selectedCard.ImageUris.Large},
	}

	generatedTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	return generatedTable
}

func initialModel() model {
	textInput := textinput.New()
	textInput.Placeholder = "Jace, Mind Sculptor"
	textInput.PromptStyle = ui.FocusedStyle
	textInput.Focus()

	textInput.CharLimit = 156
	textInput.Width = 20
	s := ui.Spinner()

	return model{
		textInput: textInput,
		mode:      Search,
		spinner:   s,
		cardList:  list.New([]list.Item{}, ui.ItemDelegate{}, 0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.mode {
			case Search:
				m.mode = Select
				m.isLoading = true
				return m, tea.Batch(m.spinner.Tick, getCardsCmd(m.textInput.Value()))
			case Select:
				m.selectedCard = (m.cardList.SelectedItem().(*ui.CardListItem)).Card
				m.infoTable = m.setInfoTable()
				m.infoTable.SetStyles(ui.TableStyle())
				m.mode = View
			}
		case "b":
			switch m.mode {
			case Select:
				m.mode = Search
			case View:
				m.mode = Select
			}
		}

	case getCardsMsg:
		m.cardList = list.New(msg.cards, ui.ItemDelegate{}, 20, 14)
		fmt.Printf("Fetched card details: %+v\n", m.cardList)
		m.isLoading = false
	}

	switch m.mode {
	case Search:
		m.textInput, cmd = m.textInput.Update(msg)
	case Select:
		m.cardList, cmd = m.cardList.Update(msg)
	case View:
		m.infoTable, cmd = m.infoTable.Update(msg)
	}
	var sCmd tea.Cmd
	m.spinner, sCmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd, sCmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.isLoading {
		return fmt.Sprintf("Loading... %s", m.spinner.View())
	}
	switch m.mode {
	case Search:
		return fmt.Sprintf(
			"Enter a card name: \n%s\n",
			m.textInput.View(),
		) + ui.HelpStyle("\n enter: choose • q/ctrl+c: quit\n")
	case Select:
		m.cardList.Title = "Select a card"
		return fmt.Sprintf(
			m.cardList.View(),
		) + ui.HelpStyle("\n enter: choose • b: back\n")
	case View:
		return fmt.Sprintf(
			m.infoTable.View(),
		) + ui.HelpStyle("\n b: back • q/ctrl+c: quit\n")
	}

	return "Unknown mode"
}

func main() {
	app := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
