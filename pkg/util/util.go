package util

import (
	"errors"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

func PromptSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	return result, err
}

func CreateFile(folder string, fileName string, content string) error {
	fileName = filepath.Join(folder, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return  err
	}

	_, err = file.WriteString(content)
	return err
}

func PromptYN(promtText string) bool {
	result, _ := PromptString(promtText + ": y/n")
	return result == "y" || result == "Y"
}

func ExecStringCommand(command string) (string, error) {
	args := strings.Fields(command)
	result, err := exec.Command(args[0], args[1:]...).Output()
	return string(result), err
}

func GetCurrentFolderName() (string, error){
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	parent := filepath.Base(pwd)

	return parent, err
}

func StripNewLine(input string) string {
	return strings.Replace(input, "\n", "", -1)
}

func HandleFatalError(err error) {
	if err != nil {
		exitError, isExitError := err.(*exec.ExitError)
		if isExitError {
			log.Print(string(exitError.Stderr))
		}
		log.Fatal(err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}