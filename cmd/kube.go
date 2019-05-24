package cmd

import (
	"errors"
	"github.com/ryannel/hippo/internal/kube"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
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
		err := kube.Deploy(envName)
		util.HandleFatalError(err)

		log.Print("Deployment Completed.")
	},
}


