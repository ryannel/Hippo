package cmd

import (
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/docker"
	componentEnum "github.com/ryannel/hippo/pkg/enum/components"
	"github.com/ryannel/hippo/pkg/enum/dockerRegistries"
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

		log.Print(componentEnum.Docker + " component has been created")
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
	},
}

func setupDocker(projectFolderPath string) {
	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)
	config := confManager.GetConfig()

	scaffold, err := scaffoldManager.New(config.ProjectName, projectFolderPath, config.Language)
	util.HandleFatalError(err)

	err = scaffold.CreateDockerFile()
	util.HandleFatalError(err)

	err = scaffold.CreateDockerIgnore()
	util.HandleFatalError(err)

	dockerRegistry, err := util.PromptSelect("Docker Registry", []string {dockerRegistries.QuayIo, "None"})
	util.HandleFatalError(err)

	if dockerRegistry == "None" {
		return
	}

	err = confManager.SetDockerRegistry(dockerRegistry)
	util.HandleFatalError(err)

	registryDomain := docker.GetRegistryDomain(dockerRegistry)

	err = confManager.SetDockerRegistryDomain(registryDomain)
	util.HandleFatalError(err)

	registryNamespace, err := util.PromptString("Docker Registry Namespace")
	util.HandleFatalError(err)

	err = confManager.SetDockerRegistryNamespace(registryNamespace)
	util.HandleFatalError(err)

	err = confManager.SetDockerRegistryUrl(registryDomain + "/" + registryNamespace)
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

func setupWizard(projectFolderPath string) {

}



