package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"screen-manager/utils"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
)

type dettachScreen struct{}

type listKeyMap struct {
	toggleStatusBar key.Binding
	toggleHelpMenu  key.Binding
	attachScreen    key.Binding
	startScreen     key.Binding
	stopScreen      key.Binding
	quit            key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S", "s"),
			key.WithHelp("S", "toggle status"),
		),
		attachScreen: key.NewBinding(
			key.WithKeys("A", "a"),
			key.WithHelp("A", "attach screen"),
		),
		startScreen: key.NewBinding(
			key.WithKeys("R", "r"),
			key.WithHelp("R", "start screen"),
		),
		stopScreen: key.NewBinding(
			key.WithKeys("D", "d"),
			key.WithHelp("D", "stop screen"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H", "h"),
			key.WithHelp("H", "toggle help"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "Q", "esc", "ctrl+c"),
			key.WithHelp("Q", "Quit"),
		),
	}
}

type model struct {
	quitting bool
	list     list.Model
	keys     *listKeyMap
}

func newModel() model {
	var (
		listKeys = newListKeyMap()
	)

	// Make initial list of items
	screenConfs := utils.GetAllConfigs()
	items := make([]list.Item, len(screenConfs))
	for i, conf := range screenConfs {
		items[i] = conf
	}

	// Setup list
	screenList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	screenList.Title = "Screens"
	screenList.Styles.Title = titleStyle
	screenList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.startScreen,
			listKeys.stopScreen,
			listKeys.toggleStatusBar,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list: screenList,
		keys: listKeys,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case dettachScreen:
		tea.ExitAltScreen()
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.startScreen):
			m.list.SelectedItem().(utils.ScreenConfig).Run()
			return m, nil

		case key.Matches(msg, m.keys.stopScreen):
			m.list.SelectedItem().(utils.ScreenConfig).Stop()
			return m, nil

		case key.Matches(msg, m.keys.attachScreen):
			tea.EnterAltScreen()
			cmd, valid := m.list.SelectedItem().(utils.ScreenConfig).Attach()
			if valid {
				command := tea.ExecProcess(cmd, func(err error) tea.Msg {
					return dettachScreen{}
				})
				return m, command
			}
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return appStyle.Render(m.list.View())
}

func main() {
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
