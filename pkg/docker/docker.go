package docker

import (
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"strings"
)

var execStringCommand = util.ExecStringCommand
var execCommandStreamingOut = util.ExecCommandStreamingOut

func Build(imageName string, commitTag string) error {
	commitTag = generateTag(imageName, commitTag)
	if commitTag != "" {
		commitTag = "-t " + commitTag + " "
	}

	command := strings.ToLower(`docker build --pull ` + commitTag + `.`)
	logger.Command("Building Docker image: " + command)
	return execCommandStreamingOut(command)
}

func Tag(sourceImage string, sourceTag string, targetImage string, targetTag string) error {
	sourceTag = generateTag(sourceImage, sourceTag)
	targetTag = generateTag(targetImage, targetTag)

	command := strings.ToLower("docker tag " + sourceTag + " " + targetTag)
	logger.Command("Tagging image : " + command)
	result, err := execStringCommand(command)
	if err != nil {
		return err
	}
	logger.Log(result)
	return err
}

func Push(registryUrl string, imageName string, tag string) error {
	tag = generateTag(imageName, tag)
	command := strings.ToLower("docker push " + registryUrl + "/" + tag)
	logger.Command("Pushing docker image: " + command)
	err := execCommandStreamingOut(command)
	return err
}

func Login(registryUrl string, username string, password string) error {
	command := "docker login -u " + username + " -p " + password + " " + registryUrl
	logger.Command("Logging into docker registry: docker login -u " + username + " -p <password> " + registryUrl)
	result, err := execStringCommand(command)
	if err != nil {
		return err
	}
	logger.Log(result)
	return err
}

func generateTag(name string, tag string) string {
	arg := name

	if tag != "" {
		arg = arg + ":" + tag
	}

	return arg
}
