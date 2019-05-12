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

	vcs, err := versionControl.New()
	util.HandleFatalError(err)

	imageName := config.ProjectName
	commit := getCommit(vcs)
	registryUrl := getRegistryUrl(config)

	err = docker.Build(registryUrl, imageName, commit)
	util.HandleFatalError(err)

	err = docker.Tag(registryUrl, commit, getTag(vcs))
	util.HandleFatalError(err)

	log.Print("Build Completed.")
}

func getCommit(vcs versionControl.VersionControl) string {
	commit, err := vcs.GetCommit()
	util.HandleFatalError(err)

	return commit
}

func getRegistryUrl(config configManager.Config) string {
	registryUrl := config.DockerRegistryUrl
	if registryUrl == "" {
		err := errors.New("`DockerRegistryUrl` not configured. Run `hippo setup docker` to enable repository handling")
		util.HandleFatalError(err)
	}

	return registryUrl
}

func getConfig() configManager.Config {
	confManager, err := configManager.New("congfig.yaml")
	util.HandleFatalError(err)

	config, err := confManager.ParseConfig()
	util.HandleFatalError(err)

	return config
}

func getTag(vcs versionControl.VersionControl) string {
	tag, err := vcs.GetBranchReplaceSlash()
	util.HandleFatalError(err)

	return tag
}


