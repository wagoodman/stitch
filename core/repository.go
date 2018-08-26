package core

import (
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Repository struct {
	Name    string
	GitURL  string `yaml:"git"`
	Path    string
	Version string
}

// NewRepository creates a new Repository populated with sane default values
func NewRepository() (obj Repository) {
	obj.Version = "master"
	return obj
}

// UnmarshalYAML parses and creates a Repository from a given user yaml string
func (repository *Repository) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type defaults Repository
	values := defaults(NewRepository())

	if err := unmarshal(&values); err != nil {
		return err
	}

	*repository = Repository(values)

	// set derivative values or overrides here here
	repository.Path = repository.DefaultRepoPath()

	// ensure that all fields are valid
	if err := repository.validate(); err != nil {
		return err
	}

	return nil
}

// todo!
func (Repository *Repository) validate() error {
	return nil
}

func (repository *Repository) DefaultRepoPath() string {
	workspaceDir := viper.GetString("workspace-path")

	elements := strings.Split(repository.GitURL, "/")
	repoSafeName := strings.TrimSuffix(elements[len(elements)-1], ".git")
	clonePath := filepath.Join(workspaceDir, repoSafeName)
	return clonePath
}

func (repository *Repository) Clone() error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", filepath.Join(homeDir, "/.ssh/id_rsa"), "")
	Check(err, "cannot use ssh keys")

	_, err = git.PlainClone(repository.Path, false, &git.CloneOptions{
		URL:      repository.GitURL,
		Progress: os.Stdout,
		Auth:     sshAuth,
	})
	return err
}
