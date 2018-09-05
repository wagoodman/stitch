package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Project struct {
	Name       string
	Compose    *Compose
	Repository Repository
	Repositories []Repository
}

func (project *Project) SafeName() string {
	cleanStr := strings.Replace(project.Name, " ", "", -1)
	cleanStr = strings.Replace(cleanStr, "-", "", -1)
	return cleanStr
}

func ReadConfig(url string) (*Project, error) {
	project := Project{}

	// create a repository for the stitch-project
	repository := NewRepository(url, "")
	if !PathExists(repository.Path) {
		err := repository.Clone()
		if err != nil {
			return nil, err
		}
	}

	// search for stitch-project file
	projectFilePath, found := FindProjectFile(repository.Path)
	if !found {
		return nil, fmt.Errorf("could not find stitch-project")
	}

	yamlString, err := ioutil.ReadFile(projectFilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot to read yaml")
	}

	err = yaml.Unmarshal(yamlString, &project)
	if err != nil {
		return nil, fmt.Errorf("cannot parse yaml")
	}

	// override any potentially configured values
	project.Repository = *repository
	project.Compose = nil

	return &project, nil
}

func (project *Project) Load() error {
	// ensure all repos have (probably) been cloned
	for _, repository := range project.Repositories {
		if !PathExists(repository.Path) {
			return fmt.Errorf("not all project repositories exist (%s)", repository.Name)
		}
	}

	// stitch the docker-compose project together
	// todo...
	project.Compose = &Compose{project.Repositories[0].Path}

	return nil
}

// DefaultProject creates a new Project populated with sane default values
func DefaultProject() (obj Project) {
	// currently no defaults are necessary
	// obj.somefield = "default value"
	return obj
}

// UnmarshalYAML parses and creates a Project from a given user yaml string
func (project *Project) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type defaults Project
	defaultValues := defaults(DefaultProject())

	if err := unmarshal(&defaultValues); err != nil {
		return err
	}

	*project = Project(defaultValues)

	// set derivative values or overrides here here
	// project.somefield = "override value"

	// ensure that all fields are valid
	if err := project.validate(); err != nil {
		return err
	}

	return nil
}

// todo!
func (project *Project) validate() error {
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
