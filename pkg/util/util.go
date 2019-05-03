package util

import (
	"errors"
	"github.com/manifoldco/promptui"
	"strconv"
)

func PromptString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	return prompt.Run()
}

func PromptInt(label string) (int, error){
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid int")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	intResult, err := strconv.Atoi(result)

	return intResult, err
}

func PromptPassword(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Mask:     '*',
	}

	return prompt.Run()
}

func Select(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	return result, err
}

