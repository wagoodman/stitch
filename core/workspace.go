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
	// todo: add url to local path lookup
	ProjectNames      map[string]string
	CurrentProjectUrl string
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
			ws.ProjectNames = make(map[string]string)
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

func LoadCurrentProject() (*Workspace, *Project, error) {
	workspace := GetWorkspace()
	project, err := ReadConfig(workspace.CurrentProjectUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load workspace: %+v\n", err)
	}

	err = project.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load project: %+v\n", err)
	}

	return workspace, project, nil
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

func (workspace *Workspace) DoesProjectExist(name, url string) bool {
	if url != "" {
		if _, ok := workspace.ProjectNames[url]; ok {
			return true
		}
	}
	if name != "" {
		for _, projName := range workspace.ProjectNames {
			if projName == name {
				return true
			}
		}
	}

	return false
}


func (workspace *Workspace) AddProject(project *Project) error {
	if workspace.DoesProjectExist("", project.Repository.GitURL) {
		return fmt.Errorf("project already exists (%s)", project.Repository.GitURL)
	}
	workspace.ProjectNames[project.Repository.GitURL] = project.Name
	return nil
}