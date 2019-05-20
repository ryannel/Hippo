package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	componentEnum "github.com/ryannel/hippo/pkg/enum/components"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	setupCmd.AddCommand(setupDockerCmd)
	setupCmd.AddCommand(setupDbCmd)
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

		log.Print(componentEnum.Docker + " component has been created")
	},
}

func setupDocker(projectFolderPath string) {
	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)

	config, err := confManager.ParseConfig()
	util.HandleFatalError(err)

	scaffold, err := scaffoldManager.New(config.ProjectName, projectFolderPath, config.Language)
	util.HandleFatalError(err)

	err = scaffold.CreateDockerFile()
	util.HandleFatalError(err)

	err = scaffold.CreateDockerIgnore()
	util.HandleFatalError(err)

	dockerRegistryUrl, err := util.PromptString("Docker Registry Url")
	util.HandleFatalError(err)

	err = confManager.SetDockerRegistryUrl(dockerRegistryUrl)
	util.HandleFatalError(err)

	dockerRegistryUser, err := util.PromptString("Docker Registry Username")
	util.HandleFatalError(err)

	err = confManager.SetDockerRegistryUser(dockerRegistryUser)
	util.HandleFatalError(err)

	dockerRegistryPassword, err := util.PromptPassword("Docker Registry Password")
	util.HandleFatalError(err)

	err = confManager.SetDockerRegistryPassword(dockerRegistryPassword)
	util.HandleFatalError(err)
}

var setupDbCmd = &cobra.Command{
	Use:   "db <kubernetes context>",
	Short: "creates a prostgresql instance in kubernetes and assigns login secrets",
	Long:  `creates a prostgresql instance in kubernetes and assigns login secrets

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a kubernetes context name argument")
		}

		if len(args) > 1 {
			return errors.New("kubernetes context names can't include spaces")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]
		setupDb(contextName)

		log.Print(componentEnum.Db + " component has been created")
	},
}

func setupDb(kubernetesContext string) {
	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)

	config, err := confManager.ParseConfig()
	util.HandleFatalError(err)

	kubernetes.New()
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

		log.Print(componentEnum.Docker + " component has been setup")
	},
}

func setupWizard(projectFolderPath string) {

}



