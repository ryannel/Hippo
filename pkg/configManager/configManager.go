package configManager

import (
	"errors"
	"github.com/ryannel/hippo/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func New(configPath string) (configManager, error) {
	exists, err := util.PathExists(configPath)
	if err != nil || !exists {
		return configManager{}, errors.New("failed to load config file")
	}

	if err != nil {
		return configManager{}, err
	}

	return configManager{configPath}, nil
}

type Config struct {
	ProjectName string `yaml:"ProjectName"`
	Language string `yaml:"Language"`
	DockerRegistryUrl string `yaml:"DockerRegistryUrl"`
}

type configManager struct {
	configPath string
}

func (manager configManager) SetProjectName(projectName string) error {
	return manager.writeToConfig("ProjectName", projectName)
}

func (manager configManager) SetLanguage(language string) error {
	return manager.writeToConfig("Language", language)
}

func (manager configManager) writeToConfig(key string, value string) error {
	file, err := os.OpenFile(manager.configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(key + ": " + value + "\n")
	if err != nil {
		return err
	}

	return file.Close()
}

func (manager configManager) ParseConfig() (Config, error) {
	var config Config

	configYaml, err := ioutil.ReadFile(manager.configPath)
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

func CreateConfigFile(path string) error  {
	configPath := filepath.Join(path, "hippo.yaml")
	_, err := os.Create(configPath)
	return err
}