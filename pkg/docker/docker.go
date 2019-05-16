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
	command := strings.ToLower("docker tag "+sourceImage+ ":" + sourceTag+ " " +targetImage+":"+ targetTag)
	log.Print("Tagging image : " + command)
	_, err := util.ExecStringCommand(command)

	return err
}

func Push(registryUrl string, name string, tag string) error {
	tag = generateTag(name, tag)
	command := "docker push " + registryUrl + "/" + tag
	log.Print("Pushing docker image: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func Login(registryUrl string, username string, password string) error {
	command := "docker login -u " + username + " -p " + password + " " + registryUrl
	log.Print("Logging into docker registry: " + command)
	_, err := util.ExecStringCommand(command)

	return err
}

func generateTag(name string, tag string) string {
	var arg string
	if name != "" {
		arg = name
	}

	if tag != "" {
		arg = arg + ":" + tag
	}

	return arg
}
