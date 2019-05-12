package cmd

import (
	"github.com/ryannel/hippo/pkg/configManager"
	componentEnum "github.com/ryannel/hippo/pkg/enum/components"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	setupCmd.AddCommand(setupDockerCmd)
	setupCmd.AddCommand(setupWizardCmd)
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup <component>",
	Short: "Creates the configuration and files needed for a component",
	Long:  `Creates the configuration and files needed for a component

Some usage examples.
`,
}

var setupDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Creates the configuration and files needed for Docker",
	Long:  `Creates the configuration and files needed for Docker

Some usage examples.
`,
	Run: func(cmd *cobra.Command, args []string) {
		projectFolderPath, err := os.Getwd()
		util.HandleFatalError(err)

		setupDocker(projectFolderPath)

		log.Print(componentEnum.Docker + " component has been created setup")
	},
}

var setupWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Launches the setup wizard",
	Long:  `Launches the setup wizard which will prompt you through the hippo setup process.

Some usage examples.
`,
	Run: func(cmd *cobra.Command, args []string) {
		projectFolderPath, err := os.Getwd()
		util.HandleFatalError(err)

		setupWizard(projectFolderPath)

		log.Print(componentEnum.Docker + " component has been created setup")
	},
}

func setupDocker(projectFolderPath string) {
	confManager, err := configManager.New("config.yaml")
	util.HandleFatalError(err)

	config, err := confManager.ParseConfig()
	util.HandleFatalError(err)

	scaffold, err := scaffoldManager.New(config.ProjectName, projectFolderPath, config.Language)
	util.HandleFatalError(err)

	err = scaffold.CreateDockerFile()
	util.HandleFatalError(err)

	err = scaffold.CreateDockerIgnore()
	util.HandleFatalError(err)
}

func setupWizard(projectFolderPath string) {

}



