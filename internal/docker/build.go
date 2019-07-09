package docker

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/versionControl"
	"log"
)

func Build() error {
	err := dockerFileExists()
	if err != nil {
		return err
	}

	config, err := configuration.New()
	if err != nil {
		return err
	}

	if config.ConfigPath == "" {
		return errors.New("no hippo.yaml found in path. Please run `hippo configure`")
	}

	registryUrl := docker.BuildReigistryUrl(config.Docker.RegistryName, config.Docker.RegistrySubDomain, config.Docker.Namespace)

	imageName := generateDockerImageName(registryUrl, config.ProjectName)

	vcs, err := versionControl.New(config.VersionControl.Provider, config.VersionControl.NameSpace, config.VersionControl.Project, config.VersionControl.Repository, config.VersionControl.Username, config.VersionControl.Password)
	if err != nil {
		return err
	}

	commitTag, err := vcs.GetCommit()
	if err != nil {
		return err
	}

	branchTag, err := vcs.GetBranchReplaceSlash()
	if err != nil {
		return err
	}

	err = docker.Build(imageName, commitTag)
	if err != nil {
		return err
	}

	if branchTag != "" {
		err = docker.Tag(imageName, commitTag, imageName, branchTag)
		if err != nil {
			return err
		}
	}

	if branchTag == "master" {
		err = docker.Tag(imageName, commitTag, imageName, "latest")
		if err != nil {
			return err
		}
	}

	log.Print("Build Completed.")
	return nil
}
