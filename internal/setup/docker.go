package setup

import (
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/enum/dockerRegistries"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
)

func SetupDocker(projectFolderPath string) error {
	confManager, err := configManager.New("hippo.yaml")
	if err != nil {
		return err
	}
	config := confManager.GetConfig()

	scaffold, err := scaffoldManager.New(config.ProjectName, projectFolderPath, config.Language)
	if err != nil {
		return err
	}

	err = scaffold.CreateDockerFile()
	if err != nil {
		return err
	}

	err = scaffold.CreateDockerIgnore()
	if err != nil {
		return err
	}

	dockerRegistry, err := util.PromptSelect("Docker Registry", []string{dockerRegistries.QuayIo, "None"})
	if err != nil {
		return err
	}

	if dockerRegistry == "None" {
		return nil
	}

	err = confManager.SetDockerRegistry(dockerRegistry)
	if err != nil {
		return err
	}

	registryDomain := docker.GetRegistryDomain(dockerRegistry)

	registryNamespace, err := util.PromptString("Docker Registry Namespace")
	if err != nil {
		return err
	}

	err = confManager.SetDockerRegistryUrl(registryDomain + "/" + registryNamespace)
	if err != nil {
		return err
	}

	dockerRegistryUser, err := util.PromptString("Docker Registry Username")
	if err != nil {
		return err
	}

	err = confManager.SetDockerRegistryUser(dockerRegistryUser)
	if err != nil {
		return err
	}

	dockerRegistryPassword, err := util.PromptPassword("Docker Registry Password")
	if err != nil {
		return err
	}

	err = confManager.SetDockerRegistryPassword(dockerRegistryPassword)
	if err != nil {
		return err
	}

	return nil
}
