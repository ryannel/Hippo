package kube

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/util"
	"os"
	"path/filepath"
)

func KubeDeploy(envName string) error {
	projectFolder, err := os.Getwd()
	if err != nil {
		return err
	}

	deployYamlPath := filepath.Join(projectFolder, "deployment_files", "deploy.yaml")

	exists, err := util.PathExists(deployYamlPath)
	if !exists || err != nil {
		return errors.New("deployment files do not exist. run `hippo setup kubernetes` to create them: " + deployYamlPath)
	}

	confManager, err := configManager.New("hippo.yaml")
	config := confManager.GetConfig()

	if len(config.KubernetesContexts[envName]) == 0 {
		return errors.New("not a valid kubernetes context. Please ensure the context name exists in hippo.yaml. Run `hippo setup kubernetes` to configure")
	}

	k8, err := kubernetes.New(envName)
	if err != nil {
		return err
	}

	return k8.Deploy(deployYamlPath)
	return nil
}