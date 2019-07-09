package setup

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
)

func Kubernetes(wd string) error {
	config, err := configuration.New()
	if err != nil {
		return err
	}

	if config.ConfigPath == "" {
		return errors.New("no hippo.yaml found in path. Please run `hippo configure`")
	}

	config, err = KubernetesConfig(config)
	if err != nil {
		return err
	}

	err =  config.SaveConfig()
	if err != nil {
		return err
	}

	scaffold, err := scaffoldManager.New(config.ProjectName, wd, config.Language)
	if err != nil {
		return err
	}

	registryUrl := docker.BuildReigistryUrl(config.Docker.RegistryName, config.Docker.RegistrySubDomain, config.Docker.Namespace)
	err = scaffold.CreateDeploymentFile(registryUrl)
	if err != nil {
		return err
	}

	return nil
}

func KubernetesConfig(config configuration.Configuration) (configuration.Configuration, error) {
	if config.Docker.RegistryDomain == "" {
		return configuration.Configuration{}, errors.New("docker registry not configured. Please run `hippo setup docker` to configure")
	}

	for util.PromptYN("Add Kubernetes Environment?") {
		key, err := util.PromptString("Context name")
		if err != nil {
			return configuration.Configuration{}, err
		}

		value, err := util.PromptString("Kubectl Context, eg: --context docker-for-desktop --namespace default ")
		if err != nil {
			return configuration.Configuration{}, err
		}

		config.KubernetesContexts[key] = value
	}

	return config, nil
}
