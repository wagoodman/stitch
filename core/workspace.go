package core

import (
	"fmt"
	"os"
	"sync"

	"encoding/gob"

	"github.com/spf13/viper"
)

var ws *Workspace
var once sync.Once

type Workspace struct {
	Projects       map[string]Project
	CurrentProject string
}

// Load decodes the project config from a Gob file
func GetWorkspace() *Workspace {
	path := viper.GetString("state-path")
	once.Do(func() {
		file, err := os.Open(path)
		if err == nil {
			// read the project config from disk
			decoder := gob.NewDecoder(file)
			err = decoder.Decode(&ws)
		} else {
			// create a new config
			ws = &Workspace{}
			ws.Projects = make(map[string]Project)
		}
		file.Close()

		// // save state on application exit
		// goodbye.Register(func(ctx context.Context, sig os.Signal) {
		// 	// fmt.Printf("2: %[1]d: %[1]s\n", sig)
		// 	// todo: currently this is on any signal, we should check explicitly for exit
		// 	SaveProjectState(path)
		// })

	})
	return ws
}

// Save encodes the project config to a Gob file
func (workspace *Workspace) Save() error {
	path := viper.GetString("state-path")

	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(workspace)
	}
	file.Close()
	return err
}

func (workspace *Workspace) DoesProjectExist(name string) bool {
	if _, ok := workspace.Projects[name]; ok {
		return true
	}
	return false
}

func (workspace *Workspace) AddProject(project Project) error {
	if workspace.DoesProjectExist(project.Name) {
		return fmt.Errorf("project '%s' already exists", project.Name)
	}
	workspace.Projects[project.Name] = project
	return nil
}

func (workspace *Workspace) GetProjects() map[string]Project {
	return workspace.Projects
}
