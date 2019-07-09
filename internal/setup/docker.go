package setup

import (
	"errors"
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

	if config.ConfigPath == "" {
		return errors.New("no hippo.yaml found in path. Please run `hippo configure`")
	}

	err = DockerFiles(config.ProjectName, projectFolderPath, config.Language)
	if err != nil {
		return err
	}

	config, err = DockerConfig(config)
	if err != nil {
		return err
	}

	err =  config.SaveConfig()
	if err != nil {
		return err
	}

	return nil
}

func DockerFiles(projectName string, projectFolderPath string, language string) error {
	scaffold, err := scaffoldManager.New(projectName, projectFolderPath, language)
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

	return nil
}

func DockerConfig(config configuration.Configuration) (configuration.Configuration, error) {
	registryName := config.Docker.RegistryName
	if registryName == "" {
		registryName, err := util.PromptSelect("Docker Registry", []string{dockerRegistries.QuayIo, dockerRegistries.Azure, "None"})
		if err != nil {
			return configuration.Configuration{}, err
		}
		config.Docker.RegistryName = registryName
		config.Docker.RegistryDomain = docker.GetRegistryDomain(registryName)
	}

	if config.Docker.RegistrySubDomain == "" {
		if config.Docker.RegistryName != dockerRegistries.QuayIo {
			registrySubDomain, err := util.PromptString("Docker Registry Subdomain")
			if err != nil {
				return configuration.Configuration{}, err
			}
			config.Docker.RegistrySubDomain = registrySubDomain
		}
	}

	if registryName == "None" {
		err := config.SaveConfig()
		if err != nil {
			return configuration.Configuration{}, err
		}
		return configuration.Configuration{}, nil
	}

	if config.Docker.RegistryRepository == "" {
		config.Docker.RegistryRepository = config.ProjectName
	}

	if config.Docker.Namespace == "" {
		registryNamespace, err := util.PromptString("Docker Registry Namespace")
		if err != nil {
			return configuration.Configuration{}, err
		}
		config.Docker.Namespace = registryNamespace
	}

	if config.Docker.RegistryUser == "" {
		registryUser, err := util.PromptString("Docker Registry Username")
		if err != nil {
			return configuration.Configuration{}, err
		}
		config.Docker.RegistryUser = registryUser
	}

	if config.Docker.RegistryPassword == "" {
		registryPassword, err := util.PromptPassword("Docker Registry Password")
		if err != nil {
			return configuration.Configuration{}, err
		}
		config.Docker.RegistryPassword = registryPassword
	}

	return config, nil
}
