package models

import (
	"fmt"
	"goggles/utils"
	"os/user"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddedConfig struct{}

func (f form) configAdded() tea.Msg {
	return f.ConvertConfig
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type form struct {
	focusIndex int
	inputs     []textinput.Model
	base       tea.Model
}

func InitForm(returnModel tea.Model) form {
	m := form{
		inputs: make([]textinput.Model, 3),
		base:   returnModel,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = focusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Identifier"
		case 2:
			t.Placeholder = "Description"
			t.CharLimit = 128
		}

		m.inputs[i] = t
	}
	return m
}

func (f form) Init() tea.Cmd {
	return textinput.Blink
}

func (f form) Validate() bool {
	for j, i := range f.inputs {
		if j == 1 && len(i.Value()) < 2 {
			return false
		}
		if i.Value() == "" {
			return false
		}
	}
	return true
}

func (f form) ConvertConfig() utils.ScreenConfig {
	username, _ := user.Current()
	return utils.ScreenConfig{
		Name: f.inputs[0].Value(),
		Id:   f.inputs[1].Value(),
		Desc: f.inputs[2].Value(),
		User: username.Username,
		Path: fmt.Sprintf("./config/%s.toml", f.inputs[1].Value()),
	}
}

func (f form) CreateFile() {
	config := f.ConvertConfig()
	utils.CreateNewConfig(config)
}

func (f form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return f, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && f.focusIndex == len(f.inputs) {
				if f.Validate() {
					f.CreateFile()
					return f.base.Update(f.configAdded)
				}
				return f, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				f.focusIndex--
			} else {
				f.focusIndex++
			}

			if f.focusIndex > len(f.inputs) {
				f.focusIndex = 0
			} else if f.focusIndex < 0 {
				f.focusIndex = len(f.inputs)
			}

			cmds := make([]tea.Cmd, len(f.inputs))
			for i := 0; i <= len(f.inputs)-1; i++ {
				if i == f.focusIndex {
					// Set focused state
					cmds[i] = f.inputs[i].Focus()
					f.inputs[i].PromptStyle = focusedStyle
					f.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				f.inputs[i].Blur()
				f.inputs[i].PromptStyle = noStyle
				f.inputs[i].TextStyle = noStyle
			}

			return f, tea.Batch(cmds...)
		}
	}

	cmd := f.updateInputs(msg)

	return f, cmd
}

func (f *form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (f form) View() string {
	var b strings.Builder

	for i, input := range f.inputs {
		b.WriteString(input.View())
		if i < len(f.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if f.focusIndex == len(f.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	style := lipgloss.NewStyle().Padding(3)
	titleStyle := lipgloss.NewStyle().PaddingBottom(1)
	return style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("New Screen Config"),
			b.String(),
		),
	)
}
