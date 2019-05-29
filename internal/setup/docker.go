package setup

import (
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/enum/dockerRegistries"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
)

func Docker(projectFolderPath string) error {
	config, err := configuration.New()
	if err != nil {
		return err
	}

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

	registryName := config.Docker.RegistryName
	if registryName == "" {
		registryName, err = util.PromptSelect("Docker Registry", []string{dockerRegistries.QuayIo, "None"})
		if err != nil {
			return err
		}
		config.Docker.RegistryName = registryName
		config.Docker.RegistryDomain = docker.GetRegistryDomain(registryName)
	}

	if registryName == "None" {
		err = config.SaveConfig()
		if err != nil {
			return err
		}
		return nil
	}

	if config.Docker.RegistryRepository == "" {
		config.Docker.RegistryRepository = config.ProjectName
	}

	if config.Docker.Namespace == "" {
		registryNamespace, err := util.PromptString("Docker Registry Namespace")
		if err != nil {
			return err
		}
		config.Docker.Namespace = registryNamespace
	}

	if config.Docker.RegistryUser == "" {
		registryUser, err := util.PromptString("Docker Registry Username")
		if err != nil {
			return err
		}
		config.Docker.RegistryUser = registryUser
	}

	if config.Docker.RegistryPassword == "" {
		registryPassword, err := util.PromptPassword("Docker Registry Password")
		if err != nil {
			return err
		}
		config.Docker.RegistryPassword = registryPassword
	}

	err = config.SaveConfig()
	if err != nil {
		return err
	}

	return nil
}
