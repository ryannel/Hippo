package setup

import (
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/enum/versionControlProviders"
	"github.com/ryannel/hippo/pkg/util"
)

func VersionControl() error {
	config, err := configuration.New()
	if err != nil {
		return err
	}

	if config.VersionControl.Provider == "" {
		vcProvider, err := util.PromptSelect("Version Control Provider", []string{versionControlProviders.Azure, versionControlProviders.Git})
		if err != nil {
			return err
		}
		config.VersionControl.Provider = vcProvider
	}

	if config.VersionControl.NameSpace == "" {
		vcNameSpace, err := util.PromptString("Version Control Namespace")
		if err != nil {
			return err
		}
		config.VersionControl.NameSpace = vcNameSpace
	}

	if config.VersionControl.Project == "" {
		vcProject, err := util.PromptString("VersionControlProject")
		if err != nil {
			return err
		}
		config.VersionControl.Project = vcProject
	}

	config.VersionControl.Repository = config.ProjectName

	return config.SaveConfig()
}
