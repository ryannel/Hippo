package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"hippo/pkg/build"
	"hippo/pkg/environment"
	"log"
	"os/exec"
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
		config, err := environment.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = build.Build(config)
		if err != nil {
			exitError, isExitError := err.(*exec.ExitError)
			if isExitError {
				log.Print(string(exitError.Stderr))
			}
			log.Fatal(err)
		}

		log.Print("Build Completed.")
	},
}




