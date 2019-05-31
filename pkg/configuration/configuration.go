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

	config, err := mergeConfigs(configPaths)
	if err != nil {
		return Configuration{}, err
	}

	return config, err
}

func newFromPath(path string) (Configuration, error) {
	configPaths, err := walkConfigsFromPath(path)
	if err != nil {
		return Configuration{}, err
	}

	config, err := mergeConfigs(configPaths)
	parentConfig = config
	return config, err
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
		baseConfig.configPath = configPath

		dirConfig, err := parseConfig(configPath)
		if err != nil {
			return Configuration{}, err
		}

		err = mergo.Merge(&baseConfig, dirConfig, mergo.WithOverride)
		if err != nil {
			return Configuration{}, err
		}

		if numConfigs == len(configPaths)-1 {
			parentConfig = baseConfig
		}

		configs := configPaths[:numConfigs-1]
		numConfigs = len(configs)
	}
	return baseConfig, nil
}

type Configuration struct {
	configPath  string
	ProjectName string `yaml:"ProjectName,omitempty"`
	Language    string `yaml:"Language,omitempty"`
	Docker      struct {
		RegistryName       string `yaml:"RegistryName,omitempty"`
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
	if config.ProjectName == parentConfig.ProjectName {
		config.ProjectName = ""
	}

	if config.Language == parentConfig.Language {
		config.Language = ""
	}

	if config.Docker.RegistryName == parentConfig.Docker.RegistryName {
		config.Docker.RegistryName = ""
	}

	if config.Docker.RegistryDomain == parentConfig.Docker.RegistryDomain {
		config.Docker.RegistryDomain = ""
	}

	if config.Docker.RegistryRepository == parentConfig.Docker.RegistryRepository {
		config.Docker.RegistryRepository = ""
	}

	if config.Docker.Namespace == parentConfig.Docker.Namespace {
		config.Docker.Namespace = ""
	}

	if config.Docker.RegistryUser == parentConfig.Docker.RegistryUser {
		config.Docker.RegistryUser = ""
	}

	if config.Docker.RegistryPassword == parentConfig.Docker.RegistryPassword {
		config.Docker.RegistryPassword = ""
	}

	if config.VersionControl.Provider == parentConfig.VersionControl.Provider {
		config.VersionControl.Provider = ""
	}

	if config.VersionControl.NameSpace == parentConfig.VersionControl.NameSpace {
		config.VersionControl.NameSpace = ""
	}

	if config.VersionControl.Project == parentConfig.VersionControl.Project {
		config.VersionControl.Project = ""
	}

	if config.VersionControl.Repository == parentConfig.VersionControl.Repository {
		config.VersionControl.Repository = ""
	}

	if config.VersionControl.Username == parentConfig.VersionControl.Username {
		config.VersionControl.Username = ""
	}

	if config.VersionControl.Password == parentConfig.VersionControl.Password {
		config.VersionControl.Password = ""
	}

	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.configPath, configYaml, 0644)
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
