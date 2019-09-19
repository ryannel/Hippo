package cmd

import (
	"errors"
	"github.com/ryannel/hippo/internal/aws"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	awsConnectCmd.AddCommand(awsConnectElasticSearchCmd)
	awsConnectCmd.AddCommand(awsConnectPostgresCmd)

	awsCmd.AddCommand(awsConnectCmd)
	awsCmd.AddCommand(awsSetContextCmd)

	rootCmd.AddCommand(awsCmd)
}

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "automates aws commands",
}

var awsConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "automates aws ssh tunnels",
}

var awsConnectElasticSearchCmd = &cobra.Command{
	Use:   "elastic <profile>",
	Short: "creates an ssh tunnel to elastic search on AWS",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("profile name can't contain spaces")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		profile := "default"
		if len(args) > 0 {
			profile = args[0]
		}

		err := aws.ConnectElasticSearch("eu-west-1", profile)
		util.HandleFatalError(err)
	},
}

var awsConnectPostgresCmd = &cobra.Command{
	Use:   "postgres <profile>",
	Short: "creates an ssh tunnel to postgres RDS instance on AWS",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("profile name can't contain spaces")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		profile := "default"
		if len(args) > 0 {
			profile = args[0]
		}

		err := aws.ConnectPostgres("eu-west-1", profile)
		util.HandleFatalError(err)
	},
}

var awsSetContextCmd = &cobra.Command{
	Use:   "context <context>",
	Short: "logs in and changes to kubectl context",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("context name can't contain spaces")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		context := "local"
		if len(args) > 0 {
			context = args[0]
		}

		err := aws.SetContext(context)
		util.HandleFatalError(err)
	},
}