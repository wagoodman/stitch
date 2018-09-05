package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Repository struct {
	Name    string
	GitURL  string `yaml:"git"`
	Path    string
	Version string
}

// NewRepository creates a new Repository populated with sane default values
func NewRepository(url, version string) *Repository {
	repository := Repository{}
	repository.GitURL = url
	repository.Path = repository.DefaultRepoPath()
	repository.Name = repository.DefaultRepoName()
	repository.Version = version
	return &repository
}

// DefaultRepository creates a new Repository populated with sane default values
func DefaultRepository() (obj Repository) {
	obj.Version = "master"
	return obj
}

// UnmarshalYAML parses and creates a Repository from a given user yaml string
func (repository *Repository) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type defaults Repository
	values := defaults(DefaultRepository())

	if err := unmarshal(&values); err != nil {
		return err
	}

	*repository = Repository(values)

	// set derivative values or overrides here here
	repository.Path = repository.DefaultRepoPath()

	if repository.Name == "" {
		repository.Name = repository.DefaultRepoName()
	}

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

func (repository *Repository) DefaultRepoName() string {
	elements := strings.Split(repository.GitURL, "/")
	return strings.TrimSuffix(elements[len(elements)-1], ".git")
}

func (repository *Repository) DefaultRepoPath() string {
	workspaceDir := viper.GetString("workspace-path")
	repoSafeName := repository.DefaultRepoName()
	clonePath := filepath.Join(workspaceDir, repoSafeName)
	return clonePath
}

func (repository *Repository) Clone() error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", filepath.Join(homeDir, "/.ssh/id_rsa"), "")
	if err != nil {
		return err
	}

	_, err = git.PlainClone(repository.Path, false, &git.CloneOptions{
		URL:      repository.GitURL,
		Progress: os.Stdout,
		Auth:     sshAuth,
	})
	return err
}

func (repository *Repository) Pull() error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", filepath.Join(homeDir, "/.ssh/id_rsa"), "")
	if err != nil {
		return err
	}

	// We instance a new repository targeting the given path (the .git folder)
	repoObj, err := git.PlainOpen(repository.Path)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	worktree, err := repoObj.Worktree()
	if err != nil {
		return err
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:     sshAuth,
	})
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) GetComposeBytes() ([]byte, error) {
	var composeBytes []byte
	var err error

	composeFile, found := repository.GetComposePath()
	if found {
		composeBytes, err = ioutil.ReadFile(composeFile)
		if err != nil {
			return nil, err
		}
	}

	return composeBytes, nil
}

func (repository *Repository) GetComposePath() (string, bool) {
	var found bool
	var composeFile string

	for _, projectFile := range []string{"docker-compose.yml", "docker-compose.yaml"} {
		files, _ := filepath.Glob(filepath.Join(repository.Path, projectFile))
		if len(files) > 0 {
			found = true
			composeFile = files[0]
			break
		}
	}
	return composeFile, found
}