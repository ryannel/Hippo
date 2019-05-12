package scaffoldManager

import (
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"strings"
)

type goLangScaffold struct {}

func (scaffold goLangScaffold) CreateProjectTemplate(projectFolderPath string) error {
	err := scaffold.createMain(projectFolderPath)
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

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`)
}

func (goLangScaffold) CreateDockerFile(folder string, projectName string) error {
	dockerFile := template.GetGoDockerFile()
	dockerFile = strings.Replace(dockerFile, "{projectname}", projectName, -1)
	return util.CreateFile(folder, "Dockerfile", dockerFile)
}

//func (goLangScaffold) createGoModule(projectFolderPath string) error {
//	command := "go mod init "
//	log.Print("Init Go Module: " + command)
//
//	_, err := util.ExecStringCommand(command)
//	return err
//}

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