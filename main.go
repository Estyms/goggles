package main

import (
	"fmt"
	"os"
	"screen-manager/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	initialized bool
	screens     list.Model
}

func (m *Model) runThread() tea.Msg {
	selectedItem := m.screens.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedScreen := selectedItem.(utils.ScreenConfig)
	selectedScreen.Run()
	m.screens.Update(nil)
	return nil
}

func (m *Model) stopThread() tea.Msg {
	selectedItem := m.screens.SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedScreen := selectedItem.(utils.ScreenConfig)
	selectedScreen.Stop()
	m.screens.Update(nil)
	m.screens.CursorDown()
	m.screens.CursorUp()
	return nil
}

func initialModel() *Model {
	return &Model{}
}

func (m *Model) Initlist(height, width int) {
	conf := utils.GetConfig("test")
	list := list.New([]list.Item{conf}, list.DefaultDelegate{}, width/4, height/2)
	list.Title = "Screens"
	list.SetShowHelp(false)
	m.screens = list
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			{
				return m, tea.Quit
			}
		case "s":
			{

				return m, m.runThread
			}
		case "d":
			{
				return m, m.stopThread
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		if !m.initialized {
			m.Initlist(msg.Height, msg.Width)
			m.screens.SetSize(msg.Width-h, msg.Height-v)
			m.initialized = true
			return m, nil
		}
		m.screens.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.screens, cmd = m.screens.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return docStyle.Render(m.screens.View())
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("A error has occured : %v", err)
		os.Exit(1)
	}
}
