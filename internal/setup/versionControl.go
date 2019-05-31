package setup

import (
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/enum/versionControlProviders"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/versionControl"
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
		vcProject, err := util.PromptString("Version Control Project")
		if err != nil {
			return err
		}
		config.VersionControl.Project = vcProject
	}

	config.VersionControl.Repository = config.ProjectName

	if config.VersionControl.Username == "" {
		vcUseranme, err := util.PromptString("Version Control Username")
		if err != nil {
			return err
		}
		config.VersionControl.Username = vcUseranme
	}

	if config.VersionControl.Password == "" {
		vcPassword, err := util.PromptPassword("Version Control Password")
		if err != nil {
			return err
		}
		config.VersionControl.Password = vcPassword
	}

	err =  config.SaveConfig()
	if err != nil {
		return err
	}

	vcs, err := versionControl.New(config.VersionControl.Provider, config.VersionControl.NameSpace, config.VersionControl.Project, config.VersionControl.Repository, config.VersionControl.Username, config.VersionControl.Password)
	if err != nil {
		return err
	}

	err = vcs.Init()
	if err != nil {
		return err
	}

	err = vcs.TrackAllFiles()
	if err != nil {
		return err
	}

	err = vcs.CreateCommit("initial commit")
	if err != nil {
		return err
	}

	err = vcs.CreateRepository()
	if err != nil {
		return err
	}

	err = vcs.SetOrigin()
	if err != nil {
		return err
	}

	return nil
}
