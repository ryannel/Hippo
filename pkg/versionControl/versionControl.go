package versionControl

import (
	"errors"
	"github.com/ryannel/hippo/pkg/enum/versionControlProviders"
	"github.com/ryannel/hippo/pkg/hostingProviders/azure"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func New(provider string, namespace string, project string, repository string, username string, password string) (VersionControl, error) {
	providerHandler, err := getProvider(provider, namespace,project, repository, username, password)
	if err != nil {
		return VersionControl{}, err
	}

	return VersionControl{provider, namespace, project, repository, username, password, providerHandler}, nil
}

func getProvider(provider string, namespace string, project string, repository string, username string, password string) (hostProvider, error) {
	var hostProvider hostProvider
	var err error

	switch provider {
	case versionControlProviders.Azure: hostProvider = &azure.Provider{namespace, project, repository, username, password}
	default: err = errors.New("Unknown provider: " + provider)
	}

	return hostProvider, err
}

type VersionControl struct {
	provider      string
	namespace     string
	project       string
	repository    string
	username      string
	password      string
	hostProvider hostProvider
}

type hostProvider interface {
	CreateRepository() error
	GetRepositoryUrl() string
}

func (vcs *VersionControl) Init() error {
	command := "git init"
	log.Print("Creating Git repo: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) Validate() error {
	gitPath := filepath.Join(".", ".git")
	exists, err := util.PathExists(gitPath)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("not a git repository. Please run `hippo setup git`")
	}

	return nil
}

func (vcs *VersionControl) GetBranch() (string, error) {
	err := vcs.Validate()
	if err != nil {
		return "", err
	}

	command := "git rev-parse --abbrev-ref HEAD"
	log.Print("Getting branch name: " + command)
	branch, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return util.StripNewLine(branch), nil
}

func (vcs *VersionControl) GetBranchReplaceSlash() (string, error) {
	branch, err := vcs.GetBranch()
	if err != nil {
		return "", err
	}

	branch = strings.Replace(branch, "/", "_", -1)

	return branch, nil
}

func (vcs *VersionControl) GetCommit() (string, error) {
	err := vcs.Validate()
	if err != nil {
		return "", err
	}

	command := "git rev-parse HEAD"
	log.Print("Getting Commit: " + command)
	commit, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return util.StripNewLine(commit), nil
}

func (vcs *VersionControl) CreateRepository() error {
	return vcs.hostProvider.CreateRepository()
}

func (vcs *VersionControl) SetOrigin() error {
	err := vcs.Validate()
	if err != nil {
		return err
	}

	command := "git remote add origin " + vcs.hostProvider.GetRepositoryUrl()
	log.Print("Adding Git Origin: " + command)
	_, err = util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) CreateCommit(message string) error {
	err := vcs.Validate()
	if err != nil {
		return err
	}

	command := `git commit -m "` + message + `"`
	log.Print("Creating Git commit: " + command)
	_, err = exec.Command("git", "commit", "-m", `"`+message+`"`).Output()
	return err
}

func (vcs *VersionControl) TrackAllFiles() error {
	err := vcs.Validate()
	if err != nil {
		return err
	}

	command := "git add ."
	log.Print("Add all files to tracking: " + command)
	_, err = util.ExecStringCommand(command)
	return err
}