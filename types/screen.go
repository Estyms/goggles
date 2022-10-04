package types

import (
	"io/fs"
	"os"

	charmfs "github.com/charmbracelet/charm/fs"
)

type Screen struct {
	id string
}

func NewScreen(name string) Screen {
	return Screen{name}
}

func (s Screen) FilterValue() string {
	return s.id
}

func (s Screen) Title() string {
	return s.Read()
}

func (s Screen) Description() string {
	return ""
}

func (s Screen) Save() {
	cfs, _ := charmfs.NewFS()
	data := []byte("Hello World")
	_ = os.WriteFile("/tmp/screens/"+s.id, data, fs.FileMode(0644))
	file, _ := os.Open("/tmp/screens/" + s.id)
	cfs.WriteFile("screen-manager/"+s.id, file)
}

func (s Screen) Read() string {
	cfs, _ := charmfs.NewFS()
	data, _ := cfs.ReadFile("screen-manager/" + s.id)
	return string(data)
}
