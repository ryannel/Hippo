package environment

import (
	"errors"
	"gopkg.in/yaml.v2"
	vcsEnum "github.com/ryannel/hippo/pkg/enum/versionControl"
	"github.com/ryannel/hippo/pkg/util"
	"io/ioutil"
	"os"
)

type EnvConfig struct {
	Project string
	VersionControl string `yaml:"VersionControl"`

	AzureOrg string `yaml:"AzureOrg"`
	AzureProject string `yaml:"AzureProject"`
	AzureUser string `yaml:"AzureUser"`
	AzureToken string `yaml:"AzureToken"`
	GitHubOrg string `yaml:"GitHubOrg"`
	GitHubToken string `yaml:"GitHubToken"`

	BuildPipeline string `yaml:"BuildPipeline"`
	DockerRegistryUrl string `yaml:"DockerRegistryUrl"`

	Environments map[string]string `yaml:"Environments"`
}

func GetConfig() (EnvConfig, error){
	var config EnvConfig

	configYaml, err := ioutil.ReadFile("hippo.yaml")
	if os.IsNotExist(err) {
		return config, errors.New("hippo.yaml config file not found. run `hippo configure` to generate one")
	} else if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		return config, err
	}

	config.Project, err = util.GetCurrentFolderName()
	if err != nil {
		return config, err
	}

	err = validateEnv(config)

	return config, err
}

func validateEnv(config EnvConfig) error {
	if config.AzureUser == "" {
		return errors.New("DockerRegistryUrl not set in hippo.yaml. Run `hippo configure` to reconfigure")
	}

	if len(config.Environments) == 0 {
		return errors.New("no kubernetes environments configured. Run `hippo configure` to reconfigure")
	}

	if config.VersionControl != vcsEnum.Azure && config.VersionControl != vcsEnum.Github && config.VersionControl != vcsEnum.None {
		return errors.New("VersionControl config (" + config.VersionControl + ") is invalid. Must be either: azure, github or none")
	}

	return nil
}

func GenerateConfig() error {
	file, err := os.Create("hippo.yaml")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Docker Registry URL", "DockerRegistryUrl")
	if err != nil {
		return err
	}

	versionControl, err := util.PromptSelect("Version Control", []string{vcsEnum.Azure, vcsEnum.Github, vcsEnum.None})

	switch versionControl {
	case vcsEnum.Azure: err = generateAzureVcsConfig(file)
	case vcsEnum.Github: err = generateGitHubConfig(file)
	}
	if err != nil {
		return err
	}

	_, err = file.WriteString("Environments:\n")
	if err != nil {
		return err
	}

	for util.PromptYN("Add Kubernetes Environment?") {
		err = writeKeyValuePromptToConfig(file, "Environment Name", "Kubectl Context, eg: --context docker-for-desktop --namespace default ", "  ")
		if err != nil {
			return err
		}
	}

	return nil
}

func generateAzureVcsConfig(file *os.File) error {
	_, err := file.WriteString("VersionControl: azure\n")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Azure Organisation", "AzureOrg")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Azure Project", "AzureProject")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Azure User", "AzureUser")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Azure Token", "AzureToken")
	if err != nil {
		return err
	}

	return nil
}

func generateGitHubConfig(file *os.File) error {
	_, err := file.WriteString("VersionControlSystem: github")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Github Organisation", "GitHubOrg")
	if err != nil {
		return err
	}

	err = writePromptToConfig(file, "Github Token", "GitHubToken")
	if err != nil {
		return err
	}

	return nil
}

