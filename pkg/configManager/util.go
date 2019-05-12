package configManager

import (
	"github.com/ryannel/hippo/pkg/util"
	"os"
)

func writeKeyValuePromptToConfig(file *os.File, keyPrompt string, valuePrompt string, indent string) error {
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

func writePromptToConfig(file *os.File, promptText string, configValue string) error {
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


