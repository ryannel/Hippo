package docker

import (
	"errors"
	"github.com/ryannel/hippo/pkg/util"
	"log"
)

func Build(registryUrl string, imageName string, commit string) error {
	if registryUrl == "" {
		return errors.New("no docker registry provided. Run `hippo setup docker` to enable repository handling")
	}

	command := `docker build --pull --shm-size 256m --memory=3g --memory-swap=-1 -t ` +registryUrl+ `/` +imageName+ `:` +commit+ ` .`

	log.Print("Building Docker image: " + command)
	_, err := util.ExecStringCommand(command)

	return err
}

func Tag(registryUrl string, commit string, tag string) error {
	command := "docker tag "+registryUrl+"/"+commit+" "+registryUrl+":"+ tag
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
