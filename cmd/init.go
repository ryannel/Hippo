package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ryannel/hippo/pkg/postgresql"
	"log"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Run: func(cmd *cobra.Command, args []string) {
		err := postgresql.CreateDb("dbname", "user",  "password")
		if err != nil {
			log.Fatal(err)
		}
	},
}