package configuration

import (
	"errors"
	"github.com/imdario/mergo"
	"github.com/ryannel/hippo/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var parentConfig Configuration

func New() (Configuration, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return Configuration{}, err
	}

	configPaths, err := walkConfigsFromPath(currentDir)
	if err != nil {
		return Configuration{}, err
	}

	return mergeConfigs(configPaths)
}

func newFromPath(path string) (Configuration, error) {
	configPaths, err := walkConfigsFromPath(path)
	if err != nil {
		return Configuration{}, err
	}

	return mergeConfigs(configPaths)
}

func walkConfigsFromPath(currentDir string) ([]string, error) {
	var configs []string

	parentDir := filepath.Dir(currentDir)

	for currentDir != parentDir {
		configPath := filepath.Join(currentDir, "hippo.yaml")
		exists, err := util.PathExists(configPath)
		if exists && err == nil {
			configs = append(configs, configPath)
		}

		currentDir = parentDir
		if parentDir == filepath.Dir(currentDir) {
			parentDir = "."
		} else {
			parentDir = filepath.Dir(currentDir)
		}
	}

	return configs, nil
}

func mergeConfigs(configPaths []string) (Configuration, error) {
	baseConfig := Configuration{}
	numConfigs := len(configPaths)
	for numConfigs > 0 {
		configPath := configPaths[numConfigs-1]
		baseConfig.ConfigPath = configPath

		dirConfig, err := parseConfig(configPath)
		if err != nil {
			return Configuration{}, err
		}

		err = mergo.Merge(&baseConfig, dirConfig, mergo.WithOverride)
		if err != nil {
			return Configuration{}, err
		}

		if numConfigs == 2 {
			parentConfig = baseConfig
		}

		configs := configPaths[:numConfigs-1]
		numConfigs = len(configs)
	}
	return baseConfig, nil
}

type Configuration struct {
	ConfigPath  string `yaml:"ConfigPath,omitempty"`
	ProjectName string `yaml:"ProjectName,omitempty"`
	Language    string `yaml:"Language,omitempty"`
	Docker      struct {
		RegistryName       string `yaml:"RegistryName,omitempty"`
		RegistrySubDomain  string `yaml:"RegistrySubDomain,omitempty"`
		RegistryDomain     string `yaml:"RegistryDomain,omitempty"`
		RegistryRepository string `yaml:"RegistryRepository,omitempty"`
		Namespace          string `yaml:"NameSpace,omitempty"`
		RegistryUser       string `yaml:"RegistryUser,omitempty"`
		RegistryPassword   string `yaml:"RegistryPassword,omitempty"`
	} `yaml:"Docker,omitempty"`
	KubernetesContexts map[string]string `yaml:"KubernetesContexts,omitempty"`
	Deployments        map[string]struct {
		KubernetesContext string   `yaml:"KubernetesContext,omitempty"`
		Run               []string `yaml:"Run,omitempty"`
	} `yaml:"Deployments,omitempty"`
	VersionControl struct {
		Provider   string `yaml:"Provider,omitempty"`
		NameSpace  string `yaml:"NameSpace,omitempty"`
		Project    string `yaml:"Project,omitempty"`
		Repository string `yaml:"Repository,omitempty"`
		Username   string `yaml:"Username,omitempty"`
		Password   string `yaml:"Password,omitempty"`
	} `yaml:"VersionControl,omitempty"`
}

func (config *Configuration) SaveConfig() error {
	currentConfig := *config

	if currentConfig.ProjectName == parentConfig.ProjectName {
		currentConfig.ProjectName = ""
	}

	if currentConfig.Language == parentConfig.Language {
		currentConfig.Language = ""
	}

	if currentConfig.Docker.RegistryName == parentConfig.Docker.RegistryName {
		currentConfig.Docker.RegistryName = ""
	}

	if currentConfig.Docker.RegistrySubDomain == parentConfig.Docker.RegistrySubDomain {
		currentConfig.Docker.RegistrySubDomain = ""
	}

	if currentConfig.Docker.RegistryDomain == parentConfig.Docker.RegistryDomain {
		currentConfig.Docker.RegistryDomain = ""
	}

	if currentConfig.Docker.RegistryRepository == parentConfig.Docker.RegistryRepository {
		currentConfig.Docker.RegistryRepository = ""
	}

	if currentConfig.Docker.Namespace == parentConfig.Docker.Namespace {
		currentConfig.Docker.Namespace = ""
	}

	if currentConfig.Docker.RegistryUser == parentConfig.Docker.RegistryUser {
		currentConfig.Docker.RegistryUser = ""
	}

	if currentConfig.Docker.RegistryPassword == parentConfig.Docker.RegistryPassword {
		currentConfig.Docker.RegistryPassword = ""
	}

	if currentConfig.VersionControl.Provider == parentConfig.VersionControl.Provider {
		currentConfig.VersionControl.Provider = ""
	}

	if currentConfig.VersionControl.NameSpace == parentConfig.VersionControl.NameSpace {
		currentConfig.VersionControl.NameSpace = ""
	}

	if currentConfig.VersionControl.Project == parentConfig.VersionControl.Project {
		currentConfig.VersionControl.Project = ""
	}

	if currentConfig.VersionControl.Repository == parentConfig.VersionControl.Repository {
		currentConfig.VersionControl.Repository = ""
	}

	if currentConfig.VersionControl.Username == parentConfig.VersionControl.Username {
		currentConfig.VersionControl.Username = ""
	}

	if currentConfig.VersionControl.Password == parentConfig.VersionControl.Password {
		currentConfig.VersionControl.Password = ""
	}

	for key, value := range currentConfig.KubernetesContexts {
		if value == parentConfig.KubernetesContexts[key] {
			delete(currentConfig.KubernetesContexts, key)
		}
	}

	configPath := currentConfig.ConfigPath
	currentConfig.ConfigPath = ""

	configYaml, err := yaml.Marshal(currentConfig)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, configYaml, 0644)
}

func parseConfig(configPath string) (Configuration, error) {
	var config Configuration

	configYaml, err := ioutil.ReadFile(configPath)
	if os.IsNotExist(err) {
		return config, errors.New("hippo.yaml config file not found. run `hippo configure` to generate one")
	} else if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		return config, err
	}

	if config.KubernetesContexts == nil {
		config.KubernetesContexts = map[string]string{}
	}

	return config, err
}

func Create(path string) (Configuration, error) {
	configPath := filepath.Join(path, "hippo.yaml")
	_, err := os.Create(configPath)
	if err != nil {
		return Configuration{}, err
	}

	return newFromPath(path)
}
