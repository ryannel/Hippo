package scaffoldManager

import (
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type goLangScaffold struct {}

func (scaffold goLangScaffold) CreateProjectTemplate(projectFolderPath string, projectName string) error {
	appFolderPath, err := scaffold.CreateAppFolder(projectFolderPath)
	if err != nil {
		return err
	}

	err = scaffold.createMain(appFolderPath)
	if err != nil {
	    return err
	}

	err = scaffold.createGoModule(appFolderPath, projectName)
	if err != nil {
		return err
	}

	return nil
}

func (goLangScaffold) createMain(folder string) error {
	return util.CreateFile(folder, "main.go", `package main

func main() {
	println("hello world")
}`)
}

func (goLangScaffold) CreateGitIgnore(folder string) error {
	return util.CreateFile(folder,".gitignore", `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

.idea/

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`)
}

func (goLangScaffold) CreateDockerFile(folder string, projectName string) error {
	dockerFile := template.GoDockerFile(projectName)
	return util.CreateFile(folder, "Dockerfile", dockerFile)
}

func (goLangScaffold) createGoModule(projectFolderPath string, projectName string) error {
	cmd := exec.Command("go", "mod", "init")
	cmd.Dir = projectFolderPath
	_, err := cmd.Output()

	if err == nil {
		log.Print("Init Go Module: go mod init")
		return nil
	}

	log.Print("GoPath not found. Creating Go module under project name: go mod init " + projectName)
	cmd = exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectFolderPath
	_, err = cmd.Output()

	return err
}

func (goLangScaffold) CreateDockerIgnore(folder string) error{
	return util.CreateFile(folder,".dockerignore", `.git
DockerFile
README.md
deployment_files
.idea
.vscode
.gitignore
azure-pipelines.yml
.editorconfig`)
}

func (scaffold goLangScaffold) CreateAppFolder(projectFolderPath string) (string, error) {
	appFolderPath := filepath.Join(projectFolderPath, "app")
	log.Print("Creating Project app folder: " + appFolderPath)

	err := os.Mkdir(appFolderPath, os.ModePerm)

	return appFolderPath, err
}