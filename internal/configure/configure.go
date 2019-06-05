package configure

import (
	"github.com/ryannel/hippo/internal/setup"
	"github.com/ryannel/hippo/pkg/configuration"
	languageEnum "github.com/ryannel/hippo/pkg/enum/languages"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"os"
	"path/filepath"
)

func Configure() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := configuration.New()
	if err != nil {
		return err
	} else if config.ConfigPath == "" {
		config, err = configuration.Create(wd)
		if err != nil {
			return err
		}
	}

	if config.ProjectName == "" {
		projectName, err := util.PromptString("ProjectName")
		if err != nil {
			return err
		}
		config.ProjectName = projectName
	}

	if config.Language == "" {
		language, err := util.PromptSelect("Project Language", []string{languageEnum.GoLang, "Bare"})
		if err != nil {
			return err
		}
		config.Language = language
	}

	dockerPath := filepath.Join(wd, "Dockerfile")
	exists, err := util.PathExists(dockerPath)
	if err != nil {
		return err
	}
	if exists {
		log.Print("Dockerfile detected, setting up docker config")
		config, err = setup.DockerConfig(config)
		if err != nil {
			return err
		}
	}

	gitPath := filepath.Join(wd, ".git")
	exists, err = util.PathExists(gitPath)
	if err != nil {
		return err
	}
	if exists {
		log.Print("git detected, setting up git config")
		config, err = setup.VersionControlConfig(config)
		if err != nil {
			return err
		}
	}

	deploymentFilesPath := filepath.Join(wd, "deployment_files")
	exists, err = util.PathExists(deploymentFilesPath)
	if err != nil {
		return err
	}
	if exists {
		log.Print("kubernetes deployment files detected, setting up kubernetes config")
		config, err = setup.KubernetesConfig(config)
		if err != nil {
			return err
		}
	}

	err =  config.SaveConfig()
	if err != nil {
		return err
	}

	return nil
}
