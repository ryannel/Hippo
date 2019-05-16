package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/environment"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func init() {
	kubeCmd.AddCommand(kubeDeployCmd)
	rootCmd.AddCommand(kubeCmd)
}

var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "automates kubernetes commands",
	Long: `automates kubernetes commands

Some usage examples.
`,
}

// TODO: Add usage examples
var kubeDeployCmd = &cobra.Command{
	Use:   "deploy <environment>",
	Short: "automates kubernetes environment deployments",
	Long:  `automates kubernetes environment deployments

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
		err := kubeDeploy(envName)
		util.HandleFatalError(err)

		log.Print("Deployment Completed.")
	},
}

func kubeDeploy(envName string) error {
	projectFolder, err := os.Getwd()
	if err != nil {
		return err
	}

	deployYamlPath := filepath.Join(projectFolder, "deployment_files", "deploy.yaml")

	exists, err := util.PathExists(deployYamlPath)
	if !exists || err != nil {
		err = errors.New("deployment files do not exist. run `hippo setup kubernetes` to create them: " + deployYamlPath)
		util.HandleFatalError(err)
	}

	config, err := configManager.GetConfig("hippo.yaml")

	k8, err := kubernetes.New(envName)
	if err != nil {
		return err
	}

	return k8.Deploy(deployYamlPath)
}


