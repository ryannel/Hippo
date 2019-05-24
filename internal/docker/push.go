package docker

import (
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/versionControl"
	"log"
)

func Push() error {
	err := dockerFileExists()
	if err != nil {
		return err
	}

	confManager, err := configManager.New("hippo.yaml")
	if err != nil {
		return err
	}
	config := confManager.GetConfig()

	var commitTag string
	var branchTag string
	vcs, err := versionControl.New()
	if err == nil {
		commitTag, _ = vcs.GetCommit()
		branchTag, _ = vcs.GetBranchReplaceSlash()
	}

	imageName := config.ProjectName

	err = docker.Login(config.Docker.RegistryUrl, config.Docker.RegistryUser, config.Docker.RegistryPassword)
	if err != nil {
		return err
	}

	err = docker.Push(config.Docker.RegistryUrl, imageName, commitTag)
	if err != nil {
		return err
	}

	err = docker.Push(config.Docker.RegistryUrl, imageName, branchTag)
	if err != nil {
		return err
	}

	if branchTag == "master" {
		err = docker.Push(config.Docker.RegistryUrl, imageName, "latest")
		if err != nil {
			return err
		}
	}

	log.Print("Docker Push Completed.")
	return nil
}
