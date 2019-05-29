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

func New() (Configuration, error) {
	configPaths, err := getConfigPaths()
	if err != nil {
		return Configuration{}, err
	}

	return mergeConfigs(configPaths)
}

func getConfigPaths() ([]string, error) {
	var configs []string

	currentDir, err := os.Getwd()
	if err != nil {
		return []string{}, err
	}
	parentDir := filepath.Dir(currentDir)

	for currentDir != parentDir {
		configPath := filepath.Join(currentDir, "hippo.yaml")
		exists, err := util.PathExists(configPath)
		if exists && err == nil {
			configs = append(configs, configPath)
		}

		currentDir = parentDir
		parentDir = filepath.Dir(currentDir)
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

		configs := configPaths[:numConfigs-1]
		numConfigs = len(configs)
	}

	return baseConfig, nil
}

type Configuration struct {
	configPath string
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
	} `yaml:"VersionControl,omitempty"`
}

func (config *Configuration) SaveConfig() error {
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

	return config, err
}

func Create(path string) (Configuration, error) {
	configPath := filepath.Join(path, "hippo.yaml")
	_, err := os.Create(configPath)
	if err != nil {
		return Configuration{}, err
	}

	return New()
}
