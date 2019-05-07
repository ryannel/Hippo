package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"hippo/pkg/util"
	"hippo/pkg/scaffold"
	"log"
)

func init() {
	rootCmd.AddCommand(scaffoldCmd)
}

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold <project name>",
	Short: "Creates the base files and repository for a new project",
	Long:  `Creates the base files and repository for a new project

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an project name argument")
		}
		if len(args) > 1 {
			return errors.New("project name can't include spaces")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		language, err := util.Select("Project Language", []string{scaffold.GoLang, "Bare"})
		if err != nil {
			log.Fatal(err)
		}

		err = scaffold.Scaffold(projectName, language)
		if err != nil {
			log.Fatal(err)
		}
	},
}



