package cmd

import (
	"errors"
	"github.com/ryannel/hippo/pkg/postgresql"
	"github.com/ryannel/hippo/pkg/template"
	"log"
	"os"

	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/docker"
	componentEnum "github.com/ryannel/hippo/pkg/enum/components"
	"github.com/ryannel/hippo/pkg/enum/dockerRegistries"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	setupCmd.AddCommand(setupDockerCmd)
	setupCmd.AddCommand(setupLocalDbCmd)
	setupCmd.AddCommand(setupWizardCmd)
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup <component>",
	Short: "Creates the configuration and files needed for a component",
	Long: `Creates the configuration and files needed for a component

Some usage examples.
`,
}

var setupDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Creates the configuration and files needed for Docker",
	Long: `Creates the configuration and files needed for Docker

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
	config := confManager.GetConfig()

	scaffold, err := scaffoldManager.New(config.ProjectName, projectFolderPath, config.Language)
	util.HandleFatalError(err)

	err = scaffold.CreateDockerFile()
	util.HandleFatalError(err)

	err = scaffold.CreateDockerIgnore()
	util.HandleFatalError(err)

	dockerRegistry, err := util.PromptSelect("Docker Registry", []string{dockerRegistries.QuayIo, "None"})
	util.HandleFatalError(err)

	if dockerRegistry == "None" {
		return
	}

	err = confManager.SetDockerRegistry(dockerRegistry)
	util.HandleFatalError(err)

	registryDomain := docker.GetRegistryDomain(dockerRegistry)

	registryNamespace, err := util.PromptString("Docker Registry Namespace")
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

var setupLocalDbCmd = &cobra.Command{
	Use:   "localdb",
	Short: "creates a prostgresql instance in kubernetes and assigns login secrets",
	Long: `creates a prostgresql instance in kubernetes and assigns login secrets

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("setup localdb takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		setupLocalDb()

		log.Print(componentEnum.Db + " component has been created")
	},
}

func setupLocalDb() {
	confManager, err := configManager.New("hippo.yaml")
	util.HandleFatalError(err)
	config := confManager.GetConfig()

	k8, err := kubernetes.New("--context docker-for-desktop --namespace default")
	util.HandleFatalError(err)

	psqlTemplate := template.PostgresDeployYaml("postgres", "postgres", "postgres")

	log.Print("Creating Postgresql container")
	err = k8.Apply(psqlTemplate)
	util.HandleFatalError(err)

	log.Print("Connecting to DB instance")
	psql, err := postgresql.New("localhost", 5432, "postgres", "postgres", "postgres")
	util.HandleFatalError(err)

	log.Print("Creating dev user: `" + config.ProjectName + "` with password `" + config.ProjectName + "`")
	err = psql.CreateUser(config.ProjectName, config.ProjectName)
	log.Print(err)

	log.Print("Creating dev db: `" + config.ProjectName + "` with owner `" + config.ProjectName + "`")
	err = psql.CreateDb(config.ProjectName, config.ProjectName)
	util.HandleFatalError(err)
}

var setupWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Launches the setup wizard",
	Long: `Launches the setup wizard which will prompt you through the hippo setup process.

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
