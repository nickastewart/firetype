package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	quote     string
	input     string
	cursorPos int
	cursor    cursor.Model
	err       error
}

func initialModel() model {
	c := cursor.New()
	c.Focus()
	return model{
		quote:     "All we have to decide is what to do with the time that is given to us. ",
		cursorPos: 0,
		cursor:    c,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return cursor.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		switch msg.String() {
		case "backspace":
			if len(m.input) != 0 {
				m.input = m.input[:len(m.input)-1]
				m.cursorPos--
			}
		default:
			if m.cursorPos < len(m.quote)-1 {
				m.input += msg.String()
				m.cursorPos++
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m model) View() string {
	runes := []rune(m.quote)
	m.cursor.SetChar(string(runes[m.cursorPos]))

	var quoteRunes []rune = []rune(m.quote)
	var displayQuote string = ""
	for i := 0; i < len(m.quote); i++ {
		if i == m.cursorPos {
			displayQuote += m.cursor.View()
		} else {
			displayQuote += string(quoteRunes[i])
		}
	}

	return fmt.Sprintf("\n\n		%s \n\n		%s", displayQuote, m.input)
}
