package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"
	"strings"
)

type Project struct {
	Name       string
	Repository Repository
	Repositories []Repository
}

func NewProject(name, url, path string) (obj Project) {
	// todo: enhance constructor
	obj.Repository = NewRepository()
	obj.Name = name
	obj.Repository.GitURL = url
	obj.Repository.Path = path

	fields := strings.Split(url, "/")
	obj.Repository.Name = strings.TrimSuffix(fields[len(fields)-1], ".git")
	return obj
}

// DefaultProject creates a new Project populated with sane default values
func DefaultProject() (obj Project) {
	// currently no defaults are necessary
	// obj.somefield = "default value"
	return obj
}

// UnmarshalYAML parses and creates a ProjectConfig from a given user yaml string
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

func (project *Project) Update() error {

	// clone the stitch-project repository (if it doesn't already exist)
	if !PathExists(project.Repository.Path) {
		fmt.Println("Cloning stitch-project repository...")

		err := project.Repository.Clone()
		Check(err, "unable to clone repository")
	} else {
		// Note: we may not want to pull the latest... this is TBD
		fmt.Printf("Repository '%s' already exists, pulling latest...\n", project.Repository.Name)
		err := project.Repository.Pull()
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("...already up to date")
		} else {
			Check(err, "unable to pull latest from repository")
		}

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

	err = yaml.Unmarshal(yamlString, &project)
	if err != nil {
		return fmt.Errorf("cannot parse yaml")
	}

	for _, repository := range project.Repositories {
		if !PathExists(repository.Path) {
			fmt.Printf("Cloning '%s' repository...\n", repository.Name)

			err := repository.Clone()
			Check(err, "unable to clone repository")

		} else {
			// Note: we may not want to pull the latest...

			// fmt.Printf("Repository '%s' already exists. pulling latest...\n", repository.Name)
			// err := repository.Pull()
			// if err == git.NoErrAlreadyUpToDate {
			// 	fmt.Println("...already up to date")
			// } else {
			// 	Check(err, "unable to pull latest from repository")
			// }

			fmt.Printf("Repository '%s' already exists. \n", repository.Name)
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
