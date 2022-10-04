package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

type ScreenConfig struct {
	Name     string
	Id       string
	Desc     string `yaml:"description"`
	User     string
	Path     string
	Commands []string
}

func (sc ScreenConfig) Title() string {
	status := map[bool]string{true: "Running", false: "Stopped"}[sc.Running()]
	return fmt.Sprintf("%s : %s", sc.Name, status)
}

func (sc ScreenConfig) Description() string {
	return fmt.Sprintf("%s", sc.Desc)
}

func (sc ScreenConfig) FilterValue() string {
	return sc.Description() + sc.Title()
}

func (sc ScreenConfig) Running() bool {
	files, err := ioutil.ReadDir("/run/screen/S-" + sc.User)
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

func GetConfig(name string) ScreenConfig {
	path := fmt.Sprintf("./config/%s", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	config := ScreenConfig{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	config.Path = path
	return config
}

func GetAllConfigs() []ScreenConfig {
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
		if strings.Contains(file.Name(), ".yml") {
			screenConfs = append(screenConfs, GetConfig(file.Name()))
		}
	}

	return screenConfs
}
