package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	languageEnum "github.com/ryannel/hippo/pkg/enum/languages"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
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

		scaffoldProject(projectName, language)
	},
}

func scaffoldProject (projectName string, language string) {
	projectFolderPath := scaffoldManager.CreateProjectFolder(projectName)

	scaffold, err := scaffoldManager.New(projectName, projectFolderPath, language)

	err = scaffold.CreateProjectTemplate()
	util.HandleFatalError(err)

	err = scaffold.CreateEditorConfig()
	util.HandleFatalError(err)

	err = scaffold.CreateReadme()
	util.HandleFatalError(err)

	err = configManager.CreateConfigFile(projectFolderPath)
	util.HandleFatalError(err)

	configPath := filepath.Join(projectFolderPath, "hippo.yaml")
	_, err = os.Create(configPath)
	util.HandleFatalError(err)

	confManager, err := configManager.New(configPath)
	util.HandleFatalError(err)

	err = confManager.SetProjectName(projectName)
	util.HandleFatalError(err)

	err = confManager.SetLanguage(language)
	util.HandleFatalError(err)

	util.HandleFatalError(err)

	log.Print("Project has been created at `./" + projectName + "`")
}



