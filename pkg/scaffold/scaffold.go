package scaffold

import (
	"errors"
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

const (
	GoLang = "Go Lang"
)

func Scaffold(projectName string, language string, dockerRegistryUrl string) error {
	projectFolder, err := createProjectFolder(projectName)
	if err != nil {
		return errors.New("project folder already exists")
	}

	err = createReadMe(projectFolder, projectName)
	if err != nil {
		return err
	}

	err = createDeploymentFile(projectFolder, projectName, dockerRegistryUrl)
	if err != nil {
		return err
	}

	err = createEditorConfig(projectFolder)
	if err != nil {
		return err
	}

	err = createDockerIgnore(projectFolder)
	if err != nil {
		return err
	}

	err = createLanguageSpecificFiles(projectFolder, language)
	if err != nil {
	    return err
	}

	return nil
}

func createProjectFolder(projectName string) (string, error) {
	err := os.Mkdir(projectName, os.ModePerm)

	return projectName, err
}

func createReadMe(projectFolder string, projectName string) error {
	return util.CreateFile(projectFolder, "README.md", "#" + projectName)
}

func createDockerIgnore(projectFolder string) error{
	return util.CreateFile(projectFolder,".dockerignore", `.git
DockerFile
README.md
deployment_files
.idea
.vscode
.gitignore
azure-pipelines.yml
.editorconfig`)
}

func createEditorConfig(projectFolder string) error {
	return util.CreateFile(projectFolder, ".editorconfig", "")
}

func createDeploymentFile(projectFolder string, projectName string, dockerRegistryUrl string) error {
	deploymentFilePath := filepath.Join(projectFolder, "deployment_files")
	err := os.Mkdir(deploymentFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	templateContent := template.GetGenericDeployYaml()
	templateContent = strings.Replace(templateContent, "{projectname}", projectName, -1)
	templateContent = strings.Replace(templateContent, "{dockerRegistryUrl}", dockerRegistryUrl, -1)

	err = util.CreateFile(deploymentFilePath, "deploy.yaml", templateContent)

	return err
}

func createLanguageSpecificFiles(projectFolder string, language string) error  {
	var err error

	switch language {
	case GoLang:  err = scaffoldGoLang(projectFolder)
	}

	return err
}
