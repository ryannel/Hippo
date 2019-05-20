package configManager

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func New(configPath string) (configManager, error) {
	config, err := parseConfig(configPath)
	if err != nil {
		return configManager{}, err
	}

	return configManager{configPath, config}, nil
}

type Config struct {
	ProjectName string `yaml:"ProjectName,omitempty"`
	Language    string `yaml:"Language,omitempty"`
	Docker      struct {
		RegistryName     string `yaml:"RegistryName,omitempty"`
		RegistryDomain   string `yaml:"RegistryDomain,omitempty"`
		Namespace        string `yaml:"NameSpace,omitempty"`
		RegistryUrl      string `yaml:"RegistryUrl,omitempty"`
		RegistryUser     string `yaml:"RegistryUser,omitempty"`
		RegistryPassword string `yaml:"RegistryPassword,omitempty"`
	} `yaml:"Docker,omitempty"`
	KubernetesContexts map[string]string `yaml:"KubernetesContexts"`
}

type configManager struct {
	configPath string
	config     Config
}

func (manager *configManager) GetConfig() Config {
	return manager.config
}

func (manager *configManager) saveConfig() error {
	configYaml, err := yaml.Marshal(manager.config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(manager.configPath, configYaml, 0644)
}

func (manager *configManager) SetProjectName(projectName string) error {
	manager.config.ProjectName = projectName
	return manager.saveConfig()
}

func (manager *configManager) SetLanguage(language string) error {
	manager.config.Language = language
	return manager.saveConfig()
}

func (manager *configManager) SetDockerRegistry(registryName string) error {
	manager.config.Docker.RegistryName = registryName
	return manager.saveConfig()
}

func (manager *configManager) SetDockerRegistryDomain(registryDomain string) error {
	manager.config.Docker.RegistryDomain = registryDomain
	return manager.saveConfig()
}

func (manager *configManager) SetDockerRegistryUser(dockerRegistryUser string) error {
	manager.config.Docker.RegistryUser = dockerRegistryUser
	return manager.saveConfig()
}

func (manager *configManager) SetDockerRegistryPassword(dockerRegistryPassword string) error {
	manager.config.Docker.RegistryPassword = dockerRegistryPassword
	return manager.saveConfig()
}

func (manager *configManager) SetDockerRegistryNamespace(namespace string) error {
	manager.config.Docker.Namespace = namespace
	return manager.saveConfig()
}

func (manager *configManager) SetDockerRegistryUrl(registryUrl string) error {
	manager.config.Docker.RegistryUrl = registryUrl
	return manager.saveConfig()
}

func parseConfig(configPath string) (Config, error) {
	var config Config

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

func Create(path string) (configManager, error) {
	configPath := filepath.Join(path, "hippo.yaml")
	_, err := os.Create(configPath)
	if err != nil {
		return configManager{}, err
	}

	return New(configPath)
}
