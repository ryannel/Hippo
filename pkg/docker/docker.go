package docker

import (
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"strings"
)

func Build(name string, tag string) error {
	tag = generateTag(name, tag)
	if tag != "" {
		tag = "-t " + tag
	}

	command := strings.ToLower(`docker build --pull --shm-size 256m --memory=3g --memory-swap=-1 ` + tag + ` .`)
	log.Print("Building Docker image: " + command)
	return util.ExecCommandStreamingOut(command)
}

func Tag(sourceImage string, sourceTag string, targetImage string, targetTag string) error {
	sourceTag = generateTag(sourceImage, sourceTag)
	targetTag = generateTag(targetImage, targetTag)

	command := strings.ToLower("docker tag " +sourceImage+ " " + targetTag)
	log.Print("Tagging image : " + command)
	_, err := util.ExecStringCommand(command)

	return err
}

func Push(registryUrl string, name string, tag string) error {
	tag = generateTag(name, tag)
	command := strings.ToLower("docker push " + registryUrl + "/" + tag)
	log.Print("Pushing docker image: " + command)
	err := util.ExecCommandStreamingOut(command)
	return err
}

func Login(registryUrl string, username string, password string) error {
	command := "docker login -u " + username + " -p " + password + " " + registryUrl
	log.Print("Logging into docker registry: docker login -u " + username + " -p <password> " + registryUrl)
	_, err := util.ExecStringCommand(command)

	return err
}

func generateTag(name string, tag string) string {
	arg := name

	if tag != "" {
		arg = arg + ":" + tag
	}

	return arg
}
