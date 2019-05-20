package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/versionControl"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	dockerCmd.AddCommand(dockerBuildCmd)
	dockerCmd.AddCommand(dockerPushCmd)
	rootCmd.AddCommand(dockerCmd)
}

// TODO: Add usage examples
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "automates docker workflow commands",
	Long:  `automates docker workflow commands

Some usage examples.
`,
}

var dockerBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds and tags the local docker image",
	Long:  `builds and tags the local docker image

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("docker dockerBuild takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		dockerBuild()
	},
}

var dockerPushCmd = &cobra.Command{
	Use:   "push",
	Short: "pushes the local docker image",
	Long:  `pushes the local docker image

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("docker push takes no arguments")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		dockerPush()
	},
}

func dockerBuild() {
	err := dockerFileExists()
	util.HandleFatalError(err)

	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)
	config := confManager.GetConfig()

	var commitTag string
	var branchTag string
	vcs, err := versionControl.New()
	if err == nil {
		commitTag, _ = vcs.GetCommit()
		branchTag, _ = vcs.GetBranchReplaceSlash()
	}

	imageName := generateDockerImageName(config.Docker.RegistryUrl, config.ProjectName)

	err = docker.Build(imageName, commitTag)
	util.HandleFatalError(err)

	if branchTag != "" {
		err = docker.Tag(imageName, commitTag, imageName, branchTag)
		util.HandleFatalError(err)
	}

	if branchTag == "master" {
		err = docker.Tag(imageName, commitTag, imageName, "latest")
		util.HandleFatalError(err)
	}

	log.Print("Build Completed.")
}

func dockerPush() {
	err := dockerFileExists()
	util.HandleFatalError(err)

	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)
	config := confManager.GetConfig()

	var commitTag string
	var branchTag string
	vcs, err := versionControl.New()
	if err == nil {
		commitTag, _ = vcs.GetCommit()
		branchTag, _ = vcs.GetBranchReplaceSlash()
	}

	imageName := config.ProjectName

	err = docker.Login(config.Docker.RegistryUrl, config.Docker.RegistryUser, config.Docker.RegistryPassword)
	util.HandleFatalError(err)

	err = docker.Push(config.Docker.RegistryUrl, imageName, commitTag)
	util.HandleFatalError(err)

	err = docker.Push(config.Docker.RegistryUrl, imageName, branchTag)
	util.HandleFatalError(err)

	if branchTag == "master" {
		err = docker.Push(config.Docker.RegistryUrl, imageName, "latest")
		util.HandleFatalError(err)
	}

	log.Print("Docker Push Completed.")
}

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

