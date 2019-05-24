package docker

import (
	"errors"
	"github.com/ryannel/hippo/pkg/util"
)

func dockerFileExists() error {
	exists, err := util.PathExists("Dockerfile")
	if !exists || err != nil {
		return errors.New("dockerfile does not exist. Please run `hippo setup docker`")
	}

	return nil
}

func generateDockerImageName(registryUrl string, projectName string) string {
	var registryArg string
	if registryUrl != "" {
		registryArg = registryUrl + "/"
	}

	return registryArg + projectName
}