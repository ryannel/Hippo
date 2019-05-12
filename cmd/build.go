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
	rootCmd.AddCommand(buildCmd)
}

// TODO: Add usage examples
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the current project's docker image",
	Long:  `builds the current project's docker image

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("build takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		build()
	},
}

func build() {
	config := getConfig()

	var commitTag string
	var branchTag string
	vcs, err := versionControl.New()
	if err == nil {
		commitTag, _ = vcs.GetCommit()
		branchTag, _ = vcs.GetBranchReplaceSlash()
	}

	registryUrl := config.DockerRegistryUrl
	imageName := config.ProjectName
	name := registryUrl+"/"+imageName

	err = docker.Build(name, commitTag)
	util.HandleFatalError(err)

	if branchTag == "master" {
		err = docker.Tag(name, commitTag, name, "latest")
		util.HandleFatalError(err)
	}

	err = docker.Tag(name, commitTag, name, branchTag)
	util.HandleFatalError(err)

	log.Print("Build Completed.")
}

func getConfig() configManager.Config {
	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)

	config, err := confManager.ParseConfig()
	util.HandleFatalError(err)

	return config
}


