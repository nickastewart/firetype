package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"regexp"
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
			if isValidCharacter(msg.String()) && m.cursorPos < len(m.quote)-1 {
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

func isValidCharacter(s string) bool {
	if len(s) > 1 {
		return false
	}
	matched, err := regexp.MatchString("[a-zA-Z0-9,.;:'\"?! ]", s)
	if err != nil {
		return false
	}
	return matched
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
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
			displayQuote += style.SetString(string(quoteRunes[i])).Render()
		}
	}
	
	return fmt.Sprintf("\n\n		"+displayQuote+" \n\n		%s", m.input)
}

type Char struct {
	Char  rune
	Style lipgloss.Style
}

func (char *Char) render() string {
	char.Style.SetString(string(char.Char))
	return char.Style.Render()
}
