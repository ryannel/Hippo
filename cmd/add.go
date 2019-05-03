package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"hippo/pkg/postgresql"
	"hippo/pkg/rabbitmq"
	"hippo/pkg/util"
	"log"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "add <environment name>",
	Short: "Adds a development environment",
	Long:  `Adds a kubernetes driven development environment`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an environment name argument")
		}
		if len(args) > 1 {
			return errors.New("no more than one environment name can be provide")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		result, err := util.Select("test", []string{"RabbitMq", "PostgreSQL", "Done"})
		if err != nil {
			log.Fatal(err)
		}

		switch result {
		case "PostgreSQL": err = postgresql.PromptPostgresSetup(envName)
		case "RabbitMq": err = rabbitmq.PromptRabbitSetup(envName)
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}



