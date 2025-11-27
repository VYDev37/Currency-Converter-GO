package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginTop(1)
)

// type struct definition
type RateMsg struct {
	Rate float64
}

func Model() ModelData { // default model
	input := make([]textinput.Model, 3)

	input[0] = textinput.New()
	input[0].Placeholder = "Base currency (Default: USD)"
	input[0].CharLimit = 3
	input[0].Width = 5
	input[0].SetValue("USD")

	// 2. Setup Input TO
	input[1] = textinput.New()
	input[1].Placeholder = "To (Default: IDR)"
	input[1].CharLimit = 3
	input[1].Width = 5
	input[1].SetValue("IDR")

	input[2] = textinput.New()
	input[2].Placeholder = "Amount"
	input[2].Focus()
	input[2].CharLimit = 15
	input[2].Width = 20

	return ModelData{
		Inputs:     input,
		Loading:    true,
		FocusIndex: 0,
	}
}

// Func converter to tea
func CallFetch(from, to string) tea.Cmd {
	return func() tea.Msg {
		rate, err := GetRate(from, to)
		if err != nil {
			return err
		}

		return RateMsg{Rate: rate}
	}
}

func (m *ModelData) Submit() (tea.Model, tea.Cmd) {
	if val, err := strconv.ParseFloat(m.Inputs[2].Value(), 64); err == nil {
		m.Amount = val
		m.Loading = true
		m.Err = nil

		return m, CallFetch(m.Inputs[0].Value(), m.Inputs[1].Value())
	}

	return m, nil
}

func (m *ModelData) UpdateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

// Tea declaration part
func (m ModelData) Init() tea.Cmd {
	return CallFetch(m.Inputs[0].Value(), m.Inputs[1].Value())
}

func (m ModelData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp, tea.KeyDown, tea.KeyEnter: // to edit data and submit
			s := msg.String()
			if s == "enter" && m.FocusIndex == 2 { // key enter
				return m.Submit()
			}

			if s == "up" { // key up
				m.FocusIndex--
				if m.FocusIndex > len(m.Inputs)-1 {
					m.FocusIndex = len(m.Inputs) - 1
				}
			} else {
				m.FocusIndex++
				if m.FocusIndex < 0 {
					m.FocusIndex = 0
				}
			}

			if val, err := strconv.ParseFloat(m.Inputs[2].Value(), 64); err == nil {
				m.Amount = val
				if m.Rate > 0 {
					m.Result = m.Amount * m.Rate
				} else {
					m.Loading = true
					return m, CallFetch(m.Inputs[0].Value(), m.Inputs[1].Value())
				}
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					cmds[i] = m.Inputs[i].Focus()
				} else {
					m.Inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		}

	case RateMsg:
		m.Loading = false
		m.Rate = msg.Rate
		m.Err = nil

		if m.Amount > 0 {
			m.Result = m.Amount * m.Rate
		}

	case error:
		m.Loading = false
		m.Err = msg
	}

	cmd = m.UpdateInputs(msg)
	return m, cmd
}

func (m ModelData) View() string { // Receiver of ModelData
	p := message.NewPrinter(language.English)

	s := fmt.Sprintf("%s\n\n", titleStyle.Render(" CURRENCY CONVERTER TUI "))
	s += fmt.Sprintf("Base: %s -> Target: %s\n", strings.ToUpper(m.Inputs[0].Value()), strings.ToUpper(m.Inputs[1].Value()))

	labels := []string{"From  ", "To    ", "Amount "}
	for i := 0; i < len(m.Inputs); i++ {
		s += resultStyle.Render(labels[i])
		s += m.Inputs[i].View()
	}
	s += "\n"
	if m.Loading {
		s += "Loading rates from API...\n"
	} else if m.Err != nil {
		s += errorStyle.Render(p.Sprintf("Error: %v", m.Err))
	} else {
		from := m.Inputs[0].Value()
		to := m.Inputs[1].Value()
		if m.Rate > 0 {
			s += p.Sprintf("\nRate: %.2f %s = %.2f %s\n", m.Amount, from, m.Result, to)
		}
		if m.Result > 0 {
			s += resultStyle.Render(p.Sprintf("Result: %.2f %s", m.Result, to))
		}
	}

	s += "\n\n(Press Enter to calculate, Esc to quit)"
	return s
}
