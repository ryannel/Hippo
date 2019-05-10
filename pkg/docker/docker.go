package docker

import (
	"hippo/pkg/util"
	"log"
)

func Build(registryUrl string, imageName string, commit string) error {
	command := `docker build --shm-size 256m --memory=3g --memory-swap=-1 -t ` +registryUrl+ `/` +imageName+ `:` +commit+ ` .`
	log.Print("Building Docker image: " + command)
	_, err := util.ExecStringCommand(command)

	return err
}

func Tag() {

}

func Push() {

}
