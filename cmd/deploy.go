package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/environment"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
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
		if err != nil {
			log.Fatal(err)
		}

		err = deploy(envName, config)
		util.HandleFatalError(err)

		log.Print("Deployment Completed.")
	},
}

func deploy(envName string, config environment.EnvConfig) error {
	projectFolder, err := os.Getwd()
	if err != nil {
		return err
	}

	deployYaml := filepath.Join(projectFolder, "deployment_files", "deploy.yaml")

	k8, err := kubernetes.New(envName, config)
	if err != nil {
		return err
	}

	return k8.Deploy(deployYaml)
}


