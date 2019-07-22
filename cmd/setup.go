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
	setupCmd.AddCommand(setupKubernetesCmd)
	setupCmd.AddCommand(setupVersionControlCmd)
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

		err = setup.Docker(projectFolderPath)
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

var setupVersionControlCmd = &cobra.Command{
	Use:   "git",
	Short: "create remote git repository along with an initial commit",
	Long: `create remote git repository along with an initial commit

Some usage examples.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("setup git takes no arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := setup.VersionControl()
		util.HandleFatalError(err)

		log.Print(componentEnum.VersionControl + " component has been created")
	},
}
