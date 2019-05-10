package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/ryannel/hippo/pkg/environment"
	"log"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "creates configuration files for the current directory",
	Long:  `creates configuration files for the current directory

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("configure takes no arguments")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := environment.GenerateConfig()
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Configuration Completed.")
	},
}



