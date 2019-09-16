package util

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/ryannel/hippo/pkg/logger"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func PromptString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	return prompt.Run()
}

func PromptInt(label string) (int, error) {
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
		Label: label,
		Mask:  '*',
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

	exists, err := PathExists(fileName)
	if exists {
		log.Print(fileName + " already exists")
	} else if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	return err
}

func PromptYN(promtText string) bool {
	result, _ := PromptString(promtText + ": Y/n")
	return result == "y" || result == "Y" || result == ""
}

func ExecStringCommand(command string) (string, error) {
	args := strings.Fields(command)
	result, err := exec.Command(args[0], args[1:]...).Output()
	return string(result), err
}

func ExecCommandStreamingOut(command string) error {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)

	stdOut, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdOut)
	for scanner.Scan() {
		message := scanner.Text()
		logger.Log(message)
	}

	scanner = bufio.NewScanner(stderr)
	for scanner.Scan() {
		errorText := scanner.Text()
		return errors.New(errorText)
	}

	return nil
}

func GetCurrentFolderName() (string, error) {
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
			logger.Error(string(exitError.Stderr))
		}
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func WaitForever() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-sigCh
}
