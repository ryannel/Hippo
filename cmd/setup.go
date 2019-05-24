package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/ryannel/hippo/internal/setup"
	componentEnum "github.com/ryannel/hippo/pkg/enum/components"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	setupCmd.AddCommand(setupDockerCmd)
	setupCmd.AddCommand(setupLocalDbCmd)
	setupCmd.AddCommand(setupKubernetesCmd)
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

		err = setup.SetupDocker(projectFolderPath)
		util.HandleFatalError(err)

		log.Print(componentEnum.Docker + " component has been created")
	},
}

var setupKubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "creates the configuration and files needed for kubernetes deployments",
	Long: `creates the configuration and files needed for kubernetes deployments

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("setup kubernetes takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		workingDirectory, err := os.Getwd()
		util.HandleFatalError(err)

		err = setup.Kubernetes(workingDirectory)
		util.HandleFatalError(err)

		log.Print(componentEnum.Kubernetes + " component has been created")
	},
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
		err := setup.SetupLocalDb()
		util.HandleFatalError(err)

		log.Print(componentEnum.Db + " component has been created")
	},
}
