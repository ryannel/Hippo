package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hippo",
	Long:  `All software has versions. This is Hippo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hippo v0.0.1 -- HEAD")
	},
}
