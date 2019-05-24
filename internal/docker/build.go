package docker

import (
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/versionControl"
	"log"
)

func Build() error {
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

	imageName := generateDockerImageName(config.Docker.RegistryUrl, config.ProjectName)

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
