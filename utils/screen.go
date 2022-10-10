package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	toml "github.com/pelletier/go-toml/v2"
)

type ScreenConfig struct {
	Name     string
	Id       string
	Desc     string `toml:"description"`
	User     string
	Path     string `toml:"-"`
	Commands []string
}

var runningDot = lipgloss.NewStyle().SetString("▶").
	Foreground(lipgloss.Color("41")).
	PaddingRight(1).
	String()
var stoppedDot = lipgloss.NewStyle().SetString("⏹").
	Foreground(lipgloss.Color("9")).
	PaddingRight(1).
	String()

func (sc ScreenConfig) Title() string {
	isRunning := sc.Running()

	status := map[bool]string{true: "Running", false: "Stopped"}[isRunning]
	statusDot := map[bool]string{true: runningDot, false: stoppedDot}[isRunning]
	return fmt.Sprintf("%s %s : %s", statusDot, sc.Name, status)
}

func (sc ScreenConfig) Description() string {
	return fmt.Sprintf("%s", sc.Desc)
}

func (sc ScreenConfig) FilterValue() string {
	return sc.Name + sc.Desc
}

func (sc ScreenConfig) Running() bool {
	path := os.Getenv("GOSCREENDIR") + "/S-" + sc.User
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, x := range files {
		if strings.Contains(x.Name(), sc.Id) {
			return true
		}
	}
	return false
}

func (sc ScreenConfig) Run() tea.Msg {
	if sc.Running() {
		return nil
	}
	command := []string{"screen", "-Smd", sc.Id, "bash", "-c", strings.Join(sc.Commands, ";")}
	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return nil
}

func (sc ScreenConfig) Stop() tea.Msg {
	if !sc.Running() {
		return nil
	}
	command := []string{"screen", "-XS", sc.Id, "quit"}
	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return nil
}

func (sc ScreenConfig) Attach() (*exec.Cmd, bool) {
	if !sc.Running() {
		return nil, false
	}
	command := []string{"screen", "-r", sc.Id}
	cmd := exec.Command(command[0], command[1:]...)
	return cmd, true
}

func initScreen() {
	cmd := exec.Command("screen", "-dm", "bash", "echo")
	_ = cmd.Run()
}

func CreateNewConfig(config ScreenConfig) {
	data, err := toml.Marshal(&config)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(config.Path, data, 0666)
	if err != nil {
		panic(err)
	}
}

func GetConfig(name string) ScreenConfig {
	path := fmt.Sprintf("./config/%s", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	config := ScreenConfig{}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	config.Path = path
	return config
}

func GetAllConfigs() []ScreenConfig {
	initScreen()
	configDir, err := os.Open("./config")
	if err != nil {
		panic(err)
	}
	files, err := configDir.Readdir(0)
	if err != nil {
		panic(err)
	}

	var screenConfs []ScreenConfig
	for _, file := range files {
		if strings.Contains(file.Name(), ".toml") {
			screenConfs = append(screenConfs, GetConfig(file.Name()))
		}
	}

	return screenConfs
}
