package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"hippo/pkg/deploy"
	"hippo/pkg/environment"
	"log"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

// TODO: Add usage examples
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploys the current project into Kubernetes",
	Long:  `deploys the current project into Kubernetes

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("an environment name must be provided")
		}

		if len(args) > 1 {
			return errors.New("environment name can't contain spaces")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		config, err := environment.GetConfig()
		if os.IsNotExist(err) {
			log.Fatal("hippo.yaml config file not found. run `hippo configure` to generate one")
		} else if err != nil {
			log.Fatal(err)
		}

		err = deploy.Deploy(envName, config)
		if err != nil {
			exitError, isExitError := err.(*exec.ExitError)
			if isExitError {
				log.Print(string(exitError.Stderr))
			}
			log.Fatal(err)
		}
	},
}




