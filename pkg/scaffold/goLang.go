package scaffold

import (
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"path/filepath"
	"strings"
)

func scaffoldGoLang(projectFolder string) error {
	err := createGitIgnore(projectFolder)
	if err != nil {
		return err
	}

	err = createMain(projectFolder)
	if err != nil {
	    return err
	}

	err = createDockerFile(projectFolder)
	if err != nil {
		return err
	}

	return nil
}

func createMain(projectFolder string) error {
	return util.CreateFile(projectFolder, "main.go", `package main

func main() {
	println("hello world")
}`)
}

func createGitIgnore(projectFolder string) error {
	return util.CreateFile(projectFolder,".gitignore", `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`)
}

func createDockerFile(projectFolder string) error {
	dockerFile := template.GetGoDockerFile()
	projectName  := filepath.Base(projectFolder)

	dockerFile = strings.Replace(dockerFile, "{projectname}", projectName, -1)

	return util.CreateFile(projectFolder, "Dockerfile", dockerFile)
}

func createGoModule(projectFolder string) error {
	command := "go mod init "
	log.Print("Init Go Module: " + command)

	_, err := util.ExecStringCommand(command)
	return err
}