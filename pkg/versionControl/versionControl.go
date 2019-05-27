package versionControl

import (
	"errors"
	"github.com/ryannel/hippo/pkg/enum/versionControlProviders"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"path/filepath"
	"strings"
)

type VersionControl struct{}

func (vcs *VersionControl) GetBranch() (string, error) {
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
	command := "git rev-parse HEAD"
	log.Print("Getting Commit: " + command)
	commit, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return util.StripNewLine(commit), nil
}

func (vcs *VersionControl) Init() error {
	command := "git init"
	log.Print("Creating Git repo: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) CreateRepository() {

}

func (vcs *VersionControl) SetOrigin(url string) error {
	command := "git remote add origin " + url
	log.Print("Adding Git Origin: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) CreateCommit(message string) error {
	command := "git commit -m " + message
	log.Print("Creating Git commit: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) TrackAllFiles() error {
	command := "git add ."
	log.Print("Git add all files to tracking: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func New() (VersionControl, error) {
	gitPath := filepath.Join(".", ".git")
	exists, err := util.PathExists(gitPath)
	if err != nil {
		return VersionControl{}, err
	}

	if !exists {
		return VersionControl{}, errors.New("not a git repository. Please run `git init`")
	}

	return VersionControl{}, nil
}

func BuildvcUrl(vcProvider string, vcNameSpace string, vcProject string, vcRepo string) string {
	var vcUrl string

	switch vcProvider {
	case versionControlProviders.Azure: vcUrl = buildAzureVcUrl(vcNameSpace, vcProject, vcRepo)
	}

	return vcUrl
}

func buildAzureVcUrl(vcNamespace string, vcProject string, vcRepository string) string {
	return "http://" + vcNamespace + ".visualstudio.com/" + vcProject + "/_git/" + vcRepository
}

