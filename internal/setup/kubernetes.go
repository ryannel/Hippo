package setup

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
)

func Kubernetes(wd string) error {
	confManager, err := configManager.New("hippo.yaml")
	if err != nil {
		return err
	}

	config := confManager.GetConfig()

	if config.Docker.RegistryUrl == "" {
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

		err = confManager.AddKubernetesContext(key, value)
		if err != nil {
			return err
		}
	}

	scaffold, err := scaffoldManager.New(config.ProjectName, wd, config.Language)
	if err != nil {
		return err
	}

	err = scaffold.CreateDeploymentFile(config.Docker.RegistryUrl)
	if err != nil {
		return err
	}

	return nil
}
