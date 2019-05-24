package cmd

import (
	"errors"
	"github.com/ryannel/hippo/internal/docker"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
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
		err := docker.Build()
		util.HandleFatalError(err)
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
		err := docker.Push()
		util.HandleFatalError(err)
	},
}

