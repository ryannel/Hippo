package cmd

import (
	"errors"
	"github.com/ryannel/hippo/internal/kube"
	componentEnum "github.com/ryannel/hippo/pkg/enum/components"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	kubeDeployCmd.AddCommand(kubeDeployProjectCmd)
	kubeDeployCmd.AddCommand(kubeDeployDashboardCmd)
	kubeDeployCmd.AddCommand(kubeDeployDb)
	kubeDeployCmd.AddCommand(kubeDeployRabbit)
	kubeDeployCmd.AddCommand(kubeDeploySqs)

	kubeCmd.AddCommand(kubeDeployCmd)
	rootCmd.AddCommand(kubeCmd)
}

var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "automates kubernetes commands",
}

var kubeDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "automates kubernetes commands",
}

// TODO: Add usage examples
var kubeDeployProjectCmd = &cobra.Command{
	Use:   "project <environment> <tag>",
	Short: "automates kubernetes environment deployments",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("takes only 1 arguments")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		envName := "local"
		if len(args) > 0 {
			envName = args[0]
		}

		err := kube.DeployProject(envName)
		util.HandleFatalError(err)

		log.Print("Deployment Completed.")
	},
}

var kubeDeployDashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Deploys local kubernetes DeployDashboard",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("hippo kube ui takes no arguments")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := kube.DeployDashboard()
		util.HandleFatalError(err)

		log.Print("Deployment Completed.")
	},
}

var kubeDeployDb = &cobra.Command{
	Use:   "postgres",
	Short: "creates a prostgresql instance in kubernetes and assigns login secrets",
	Long: `creates a prostgresql instance in kubernetes and assigns login secrets

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("setup postgres takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := kube.DeployPostgres()
		util.HandleFatalError(err)

		log.Print(componentEnum.Db + " component has been created")
	},
}

var kubeDeployRabbit = &cobra.Command{
	Use:   "rabbit",
	Short: "creates a rabbit mq instance in kubernetes and assigns login secrets",
	Long: `creates a rabbit mq instance in kubernetes and assigns login secrets

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("setup localrabbit takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := kube.DeployRabbit()
		util.HandleFatalError(err)

		log.Print(componentEnum.Rabbit + " component has been created")
	},
}

var kubeDeploySqs = &cobra.Command{
	Use:   "sqs",
	Short: "creates a AWS SQS instance in kubernetes",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("setup sqs takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := kube.DeployAwsSqs()
		util.HandleFatalError(err)

		log.Print("SQS component has been created")
	},
}

