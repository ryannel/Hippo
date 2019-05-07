package environment

import (
	"errors"
	"gopkg.in/yaml.v2"
	"hippo/pkg/util"
	"io/ioutil"
	"os"
)

type EnvConfig struct {
	Test string `yaml:"test"`
	Environments map[string]string `yaml:"Environments"`
}

func GetConfig() (EnvConfig, error){
	var config EnvConfig

	configYaml, err := ioutil.ReadFile("hippo.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configYaml, &config)
	return config, err
}

func validateEnv() error {
	_, err := GetConfig()
	if os.IsNotExist(err) {
		return errors.New("hippo.yaml config file not found")
	} else if err != nil {
		return err
	}

	return nil
}

func GenerateConfig() error {
	file, err := os.Create("hippo.yaml")
	if err != nil {
		return err
	}

	err = writePromptToFile(file, "Azure User", "AzureUser")
	if err != nil {
		return err
	}

	err = writePromptToFile(file, "Azure Organisation", "AzureOrg")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Environments:\n")
	if err != nil {
		return err
	}

	for util.PromptYN("Add Kubernetes Environment?") {
		err = writeKeyValuePromptToFile(file, "Environment Name", "Kubectl Context, eg: --context docker-for-desktop --namespace local ", "  ")
		if err != nil {
			return err
		}
	}

	// TODO: Handle multiple environments
	return nil
}



func writeKeyValuePromptToFile(file *os.File, keyPrompt string, valuePrompt string, indent string) error {
	key, err := util.PromptString(keyPrompt)
	if err != nil {
		return err
	}

	value, err := util.PromptString(valuePrompt)
	if err != nil {
		return err
	}

	_, err = file.WriteString(indent + key + ": " + value + "\n")
	if err != nil {
		return err
	}
	return nil

}

func writePromptToFile(file *os.File, promptText string, configValue string) error {
	response, err := util.PromptString(promptText)
	if err != nil {
		return err
	}

	_, err = file.WriteString(configValue + ": " + response + "\n")
	if err != nil {
		return err
	}
	return nil
}
