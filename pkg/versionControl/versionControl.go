package versionControl

import (
	"hippo/pkg/environment"
	"hippo/pkg/util"
	"log"
	"strings"
)

type VersionControl struct{
	config environment.EnvConfig
}

func (vcs *VersionControl) GetBranch() (string, error) {
	command := "git rev-parse --abbrev-ref HEAD"
	log.Print("Getting branch name: " + command)
	return util.ExecStringCommand(command)
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
	return util.ExecStringCommand(command)
}

func New(config environment.EnvConfig) VersionControl {
	return VersionControl{config}
}

