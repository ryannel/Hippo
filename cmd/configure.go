package cmd

import (
	"errors"
	"github.com/ryannel/hippo/internal/configure"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "creates configuration files for the current project",
	Long:  `creates configuration files for the current project`,
	Example: `hippo configure`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("configure takes no arguments")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := configure.Configure()
		if err != nil {
			log.Fatal(err)
		}
	},
}



