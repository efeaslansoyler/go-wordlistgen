package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	inputFirstName = iota
	inputLastName
	inputBirthday
	inputRelatedWords
	inputMinLength
	inputMaxLength
	numInputs
)

var (
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	faintPromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))

	formStyle = lipgloss.NewStyle().
			Padding(1, 3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63"))

	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	errorInputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Underline(true)
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	errMsg     string
	done       bool
	width      int
}

type inputConfig struct {
	placeholder string
	focused     bool
}

var inputConfigs = []inputConfig{
	{placeholder: "firstname (and secondname if you want)", focused: true},
	{placeholder: "lastname", focused: false},
	{placeholder: "birthday (optional, DD/MM/YYYY or similar, use / to separate)", focused: false},
	{placeholder: "related words that you want to add (optional, separate with , if you enter more than one)", focused: false},
	{placeholder: "min password length (optional, default 6)", focused: false},
	{placeholder: "max password length (optional, default 12)", focused: false},
}

func NewModel() *model {
	m := &model{
		inputs: make([]textinput.Model, numInputs),
	}

	for i, config := range inputConfigs {
		t := textinput.New()
		t.PromptStyle = focusedStyle
		t.PlaceholderStyle = placeholderStyle
		t.Placeholder = config.placeholder
		if config.focused {
			t.Focus()
		}
		m.inputs[i] = t
	}
	return m
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		inputWidth := m.width / 2
		if inputWidth < 20 {
			inputWidth = 20
		}
		if inputWidth > 80 {
			inputWidth = 80
		}
		for i := range m.inputs {
			m.inputs[i].Width = inputWidth
		}
		return m, nil
	case tea.KeyMsg:
		s := msg.String()
		if m.done {
			switch s {
			case "enter":
				// TODO: Add password generation logic here
				return m, tea.Quit
			case "b", "backspace":
				m.done = false
				m.focusIndex = 0
				return m, nil
			}
			return m, nil
		}
		switch s {
		case "ctrl+c", "q", "ctrl+d", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				if idx, err := validateInputs(m.inputs); err != nil {
					m.errMsg = err.Error()
					m.focusIndex = idx
					cmds := make([]tea.Cmd, len(m.inputs))
					for i := range m.inputs {
						if i == m.focusIndex {
							cmds[i] = m.inputs[i].Focus()
						} else {
							m.inputs[i].Blur()
						}
					}
					return m, tea.Batch(cmds...)
				}
				m.done = true
				return m, nil
			} else {
				m.errMsg = ""
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		case "r":
			for i := range m.inputs {
				if m.inputs[i].Value() != "" {
					m.inputs[i].Reset()
				}
			}
			m.focusIndex = 0

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			m.errMsg = ""
			return m, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func isValidBirthday(birthday string) bool {
	for _, char := range birthday {
		if !(char == '/' || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

func validateInputs(inputs []textinput.Model) (int, error) {
	if strings.TrimSpace(inputs[inputFirstName].Value()) == "" {
		return 0, fmt.Errorf("firstname cannot be empty")
	}

	if strings.TrimSpace(inputs[inputLastName].Value()) == "" {
		return 1, fmt.Errorf("lastname cannot be empty")
	}

	birthday := strings.TrimSpace(inputs[inputBirthday].Value())
	if birthday != "" && !strings.Contains(birthday, "/") {
		return 2, fmt.Errorf("birthday must be seperated with /(format doesnt matter)")
	}
	if birthday != "" && !isValidBirthday(birthday) {
		return 2, fmt.Errorf("birthday can only contain numbers and /")
	}

	relatedWords := strings.TrimSpace(inputs[inputRelatedWords].Value())
	if relatedWords != "" {
		words := strings.Fields(relatedWords)
		if len(words) > 1 && !strings.Contains(relatedWords, ",") {
			return 3, fmt.Errorf("related words must be seperated with ,")
		}
	}

	minLength := strings.TrimSpace(inputs[inputMinLength].Value())
	if minLength != "" {
		minLen, err := strconv.Atoi(minLength)
		if err != nil || minLen < 1 {
			return 4, fmt.Errorf("min password length must be a positive number")
		}
	}

	maxLength := strings.TrimSpace(inputs[inputMaxLength].Value())
	if maxLength != "" {
		maxLen, err := strconv.Atoi(maxLength)
		if err != nil || maxLen < 1 {
			return 5, fmt.Errorf("max password length must be a positive number")
		}
	}

	if minLength != "" && maxLength != "" {
		minLen, _ := strconv.Atoi(minLength)
		maxLen, _ := strconv.Atoi(maxLength)
		if minLen > maxLen {
			return 4, fmt.Errorf("min password length cannot be greater than max password length")
		}
	}
	return -1, nil
}

func (m *model) setInputStyles() {
	for i := range m.inputs {
		// Reset to default (no style)
		m.inputs[i].PromptStyle = lipgloss.NewStyle()
		m.inputs[i].TextStyle = lipgloss.NewStyle()
		m.inputs[i].PlaceholderStyle = placeholderStyle

		if i == m.focusIndex && m.errMsg != "" {
			m.inputs[i].PromptStyle = errorInputStyle
			m.inputs[i].TextStyle = errorInputStyle
		} else if i == m.focusIndex {
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else if i != m.focusIndex && m.inputs[i].Value() != "" {
			m.inputs[i].PromptStyle = faintPromptStyle
		}
	}
}

func (m *model) View() string {
	localFormStyle := formStyle.Width(m.width/2 + 6)
	if m.done {
		minLength := 6
		maxLength := 12
		if s := strings.TrimSpace(m.inputs[4].Value()); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				minLength = v
			}
		}
		if s := strings.TrimSpace(m.inputs[5].Value()); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				maxLength = v
			}
		}
		summary := fmt.Sprintf(
			"Form Submitted!\n\nFirstname: %s\nLastname: %s\nBirthday: %s\nRelated words: %s\nMin password length: %d\nMax password length: %d\n\nPress enter to generate password.\nPress b to go back and edit.",
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			m.inputs[2].Value(),
			m.inputs[3].Value(),
			minLength,
			maxLength,
		)
		return localFormStyle.Render(summary)
	}
	var b strings.Builder

	m.setInputStyles()

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	if m.errMsg != "" {
		b.WriteString(errorStyle.Render("\n" + m.errMsg + "\n"))
	}

	help := placeholderStyle.Render(
		"\n(tab/shift+tab to move, enter to submit, r to clear, esc to quit)\n",
	)

	return localFormStyle.Render(b.String() + help)
}

func Start() {
	if _, err := tea.NewProgram(NewModel()).Run(); err != nil {
		fmt.Printf("could not start program: %v", err)
		os.Exit(1)
	}
}
