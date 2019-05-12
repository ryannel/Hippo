package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/docker"
	"github.com/ryannel/hippo/pkg/environment"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/versionControl"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

// TODO: Add usage examples
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "pushes the latest committed docker image",
	Long:  `pushes the latest committed docker image

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("push takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := environment.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = push(config)
		util.HandleFatalError(err)
		log.Print("Build Completed.")
	},
}

func push(config environment.EnvConfig) error {
	registryUrl := config.DockerRegistryUrl
	imageName := config.ProjectName

	vcs, err := versionControl.New()
	if err != nil {
		return err
	}

	commit, err := vcs.GetCommit()
	if err != nil {
		return err
	}

	tag, err := vcs.GetBranchReplaceSlash()
	if err != nil {
		return err
	}

	return docker.Push(registryUrl, imageName, tag, commit)
}



