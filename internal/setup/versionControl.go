package setup

import (
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/enum/versionControlProviders"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/versionControl"
)

func VersionControl() error {
	confManager, err := configManager.New("hippo.yaml")
	if err != nil {
		return err
	}

	config := confManager.GetConfig()

	vcProvider, err := util.PromptSelect("Version Control Provider", []string{versionControlProviders.Azure, versionControlProviders.Git})
	if err != nil {
		return err
	}
	err = confManager.SetVersionControlProvider(vcProvider)
	if err != nil {
		return err
	}

	vcNameSpace, err := util.PromptString("Version Control Namespace")
	if err != nil {
		return err
	}
	err = confManager.SetVersionControlNamespace(vcNameSpace)
	if err != nil {
	    return err
	}

	vcProject, err := util.PromptString("VersionControlProject")
	if err != nil {
	    return err
	}
	err = confManager.SetVersionControlProject(vcProject)
	if err != nil {
	    return err
	}

	vcRepo := config.ProjectName
	err = confManager.SetVersionControlRepositoryName(vcRepo)
	if err != nil {
	    return err
	}
	
	vcUrl := versionControl.BuildvcUrl(vcProvider, vcNameSpace, vcProject, vcRepo)
	err = confManager.SetVersionControlUrl(vcUrl)
	if err != nil {
	    return err
	}
}
