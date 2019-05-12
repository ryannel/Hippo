package scaffoldManager

import (
	"errors"
	languageEnum "github.com/ryannel/hippo/pkg/enum/languages"
	"github.com/ryannel/hippo/pkg/util"
	"os"
	"path/filepath"
)

func New(projectName string, projectFolderPath string, language string) (scaffold, error) {
	exists, err := util.PathExists(projectFolderPath)
	if err != nil || !exists {
		return scaffold{}, errors.New("unable to access project folder: " + projectFolderPath)
	}

	languageScaffold, err :=  selectLanguageScaffold(language)
	if err != nil {
		return scaffold{}, err
	}

	return scaffold{
		projectFolderPath: projectFolderPath,
		projectName:       projectName,
		languageScaffold:  languageScaffold,
	}, nil
}

type languageScaffold interface {
	CreateProjectTemplate(projectFolderPath string) error
	CreateGitIgnore(folder string) error
	CreateDockerFile(folder string, projectName string) error
	CreateDockerIgnore(folder string) error
}

var _ languageScaffold = goLangScaffold{}

type scaffold struct {
	projectFolderPath string
	projectName       string
	language          string
	languageScaffold  languageScaffold
}

func (scaffold scaffold) CreateProjectTemplate() error {
	return scaffold.languageScaffold.CreateProjectTemplate(scaffold.projectFolderPath)
}

func (scaffold scaffold) CreateGitIgnore() error {
	return scaffold.languageScaffold.CreateGitIgnore(scaffold.projectFolderPath)
}

func (scaffold scaffold) CreateReadme() error {
	return util.CreateFile(scaffold.projectFolderPath, "README.md", "#" + scaffold.projectName)
}

func (scaffold scaffold) CreateEditorConfig() error {
	return util.CreateFile(scaffold.projectFolderPath, ".editorconfig", "")
}

func (scaffold scaffold) CreateDockerFile() error {
	return scaffold.languageScaffold.CreateDockerFile(scaffold.projectFolderPath, scaffold.projectName)
}

func (scaffold scaffold) CreateDockerIgnore() error {
	return scaffold.languageScaffold.CreateDockerIgnore(scaffold.projectFolderPath)
}

//func createDeploymentFiles(projectFolderPath string, projectName string, dockerRegistryUrl string) error {
//	deploymentFilePath := filepath.Join(projectFolderPath, "deployment_files")
//	err := os.Mkdir(deploymentFilePath, os.ModePerm)
//	if err != nil {
//		return err
//	}
//
//	templateContent := template.GetGenericDeployYaml()
//	templateContent = strings.Replace(templateContent, "{projectname}", projectName, -1)
//	templateContent = strings.Replace(templateContent, "{dockerRegistryUrl}", dockerRegistryUrl, -1)
//
//	err = util.CreateFile(deploymentFilePath, "deploy.yaml", templateContent)
//
//	return err
//}

func selectLanguageScaffold(language string) (languageScaffold, error) {
	var scaffold languageScaffold

	switch language {
	case languageEnum.GoLang:  scaffold = goLangScaffold{}
	default: return goLangScaffold{}, errors.New("Selected Language not supported, please see configuration options in README.md: " + language)
	}

	return scaffold, nil
}

func CreateProjectFolder(projectName string) string {
	workingDirectory, err := os.Getwd()
	util.HandleFatalError(err)

	projectFolderPath := filepath.Join(workingDirectory, projectName)
	err = os.Mkdir(projectFolderPath, os.ModePerm)
	util.HandleFatalError(err)

	return projectFolderPath
}
