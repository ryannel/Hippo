package versionControl

import (
	"errors"
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

