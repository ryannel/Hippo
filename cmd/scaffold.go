package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/scaffold"
	"log"
	"strings"
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

		if strings.ToLower(args[0])[0] != args[0][0] {
			return errors.New("project name must be lower case")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		language, err := util.PromptSelect("Project Language", []string{scaffold.GoLang, "Bare"})
		if err != nil {
			log.Fatal(err)
		}

		dockerRegistryUrl, err := util.PromptString("Docker Registry URL")
		if err != nil {
			log.Fatal(err)
		}

		err = scaffold.Scaffold(projectName, language, dockerRegistryUrl)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("Project has been created at `./" + projectName + "`")
	},
}



