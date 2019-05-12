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

func Tag(sourceImage string, targetImage string, tag string) error {
	command := strings.ToLower("docker tag "+sourceImage+ " " +targetImage+":"+ tag)
	log.Print("Tagging image with tag (" + tag + "): " + command)
	_, err := util.ExecStringCommand(command)

	return err
}

func Push(registryUrl string, imageName string, tag string, commit string) error {
	command := "docker push " + registryUrl + "/" + imageName + ":" + commit
	log.Print("Pushing docker commit (" + commit + "): " + command)
	_, err := util.ExecStringCommand(command)
	if err != nil {
		return err
	}

	command = "docker push " + registryUrl + "/" + imageName + ":" + commit
	log.Print("Pushing docker tag (" + commit + "): " + tag)
	_, err = util.ExecStringCommand(command)

	return err
}

func Login(username string, password string, registryUrl string) error {
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
