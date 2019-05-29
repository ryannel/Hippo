package docker

import (
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/versionControl"
	"log"
)

func Push() error {
	err := dockerFileExists()
	if err != nil {
		return err
	}

	config, err := configuration.New()
	if err != nil {
		return err
	}

	var commitTag string
	var branchTag string
	vcs, err := versionControl.New()
	if err == nil {
		commitTag, _ = vcs.GetCommit()
		branchTag, _ = vcs.GetBranchReplaceSlash()
	}

	imageName := config.ProjectName

	registryUrl := docker.BuildReigistryUrl(config.Docker.RegistryName, config.Docker.RegistryDomain, config.Docker.Namespace, config.Docker.RegistryRepository)

	err = docker.Login(registryUrl, config.Docker.RegistryUser, config.Docker.RegistryPassword)
	if err != nil {
		return err
	}

	err = docker.Push(registryUrl, imageName, commitTag)
	if err != nil {
		return err
	}

	err = docker.Push(registryUrl, imageName, branchTag)
	if err != nil {
		return err
	}

	if branchTag == "master" {
		err = docker.Push(registryUrl, imageName, "latest")
		if err != nil {
			return err
		}
	}

	log.Print("Docker Push Completed.")
	return nil
}
