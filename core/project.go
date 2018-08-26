package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Project struct {
	Name       string
	Repository Repository
	Config     *ProjectConfig
}

type ProjectConfig struct {
	Name         string
	Repositories []Repository
}

func NewProject(name, url, path string) (obj Project) {
	pc := NewProjectConfig()
	obj.Config = &pc
	// todo: enhance constructor
	obj.Repository = NewRepository()
	obj.Name = name
	obj.Repository.GitURL = url
	obj.Repository.Path = path
	return obj
}

// NewProjectConfig creates a new ProjectConfig populated with sane default values
func NewProjectConfig() (obj ProjectConfig) {
	// currently no defaults are necessary
	// obj.somefield = "default value"
	return obj
}

// UnmarshalYAML parses and creates a ProjectConfig from a given user yaml string
func (projectConfig *ProjectConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type defaults ProjectConfig
	defaultValues := defaults(NewProjectConfig())

	if err := unmarshal(&defaultValues); err != nil {
		return err
	}

	*projectConfig = ProjectConfig(defaultValues)

	// set derivative values or overrides here here
	// projectConfig.somefield = "override value"

	// ensure that all fields are valid
	if err := projectConfig.validate(); err != nil {
		return err
	}

	return nil
}

// todo!
func (projectConfig *ProjectConfig) validate() error {
	return nil
}

// todo!
func (project *Project) validate() error {
	return nil
}

func (project *Project) Update() error {

	// clone the stitch-project repository (if it doesn't already exist)
	if !PathExists(project.Repository.Path) {
		fmt.Println("Cloning stitch-project repository...")

		err := project.Repository.Clone()
		Check(err, "unable to clone repository")
	} else {
		fmt.Println("Repository already exists")
	}

	// search for stitch-project file
	projectFilePath, found := FindProjectFile(project.Repository.Path)
	if !found {
		return fmt.Errorf("unable to add project: could not find stitch-project file in given repo")
	}

	yamlString, err := ioutil.ReadFile(projectFilePath)
	if err != nil {
		return fmt.Errorf("cannot to read yaml")
	}

	*project.Config = NewProjectConfig()
	err = yaml.Unmarshal(yamlString, project.Config)
	if err != nil {
		return fmt.Errorf("cannot parse yaml")
	}

	for _, repository := range project.Config.Repositories {
		if !PathExists(repository.Path) {
			fmt.Printf("Cloning '%s' repository...\n", repository.Name)

			err := repository.Clone()
			Check(err, "unable to clone repository")

		} else {
			fmt.Printf("Repository '%s' already exists\n", repository.Name)
		}
	}

	return nil
}

// FindProjectFile returns the path to a valid stitch-project file given a local repo path to start searching
func FindProjectFile(repoPath string) (string, bool) {
	var found bool
	var projectPath string

	for _, projectFile := range []string{".stitch-project.yaml", "stitch-project.yaml"} {
		files, err := filepath.Glob(filepath.Join(repoPath, projectFile))
		Check(err, "cannot access stitch-project path")
		if len(files) > 0 {
			found = true
			projectPath = files[0]
			break
		}
	}

	return projectPath, found
}
