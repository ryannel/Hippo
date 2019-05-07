package deploy

import (
	"hippo/pkg/environment"
	"hippo/pkg/kubernetes"
	"os"
	"path/filepath"
)

func Deploy(envName string, config environment.EnvConfig) error {
	projectFolder, err := os.Getwd()
	if err != nil {
		return err
	}

	deployYaml := filepath.Join(projectFolder, "deployment_files", "deploy.yaml")

	k8, err := kubernetes.New(envName, config)
	if err != nil {
		return err
	}

	return k8.Deploy(deployYaml)
}
