package cmd

import (
	"errors"
	"github.com/ryannel/hippo/internal/scaffold"
	languageEnum "github.com/ryannel/hippo/pkg/enum/languages"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
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

		language, err := util.PromptSelect("Project Language", []string{languageEnum.GoLang, "Bare"})
		util.HandleFatalError(err)

		err = scaffold.ScaffoldProject(projectName, language)
		util.HandleFatalError(err)
	},
}





