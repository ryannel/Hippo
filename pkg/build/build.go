package build

import (
	"hippo/pkg/docker"
	"hippo/pkg/environment"
	"hippo/pkg/util"
	"hippo/pkg/versionControl"
	"log"
)

func Build(config environment.EnvConfig) error {
	registryUrl := config.DockerRegistryUrl
	imageName := config.Project

	vcs := versionControl.New(config)
	commit, err := vcs.GetCommit()

	err = docker.Build(registryUrl, imageName, commit)
	if err != nil {
		return err
	}


	dockerTag, err := vcs.GetBranchReplaceSlash()
	if err != nil {
		return err
	}

	command := "docker tag "+registryUrl+"/"+commit+" "+registryUrl+":"+dockerTag
	log.Print("Tagging image with tag (" + dockerTag + "): " + command)
	_, err = util.ExecStringCommand(command)

	return err
}