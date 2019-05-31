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

	if config.Docker.RegistryDomain == "" {
		return errors.New("docker registry not configured. Please run `hippo setup docker` to configure")
	}

	for util.PromptYN("Add Kubernetes Environment?") {
		key, err := util.PromptString("Context name")
		if err != nil {
			return err
		}

		value, err := util.PromptString("Kubectl Context, eg: --context docker-for-desktop --namespace default ")
		if err != nil {
			return err
		}

		config.KubernetesContexts[key] = value
		err = config.SaveConfig()
		if err != nil {
			return err
		}
	}

	scaffold, err := scaffoldManager.New(config.ProjectName, wd, config.Language)
	if err != nil {
		return err
	}

	registryUrl := docker.BuildReigistryUrl(config.Docker.RegistryName, config.Docker.Namespace)
	err = scaffold.CreateDeploymentFile(registryUrl)
	if err != nil {
		return err
	}

	return nil
}
