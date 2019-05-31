package scaffoldManager

import (
	"errors"
	languageEnum "github.com/ryannel/hippo/pkg/enum/languages"
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"os"
	"path/filepath"
)

func New(projectName string, projectFolderPath string, language string) (scaffold, error) {
	exists, err := util.PathExists(projectFolderPath)
	if err != nil || !exists {
		return scaffold{}, errors.New("unable to access project folder: " + projectFolderPath)
	}

	languageScaffold, err := getLanguageScaffold(language)
	if err != nil {
		return scaffold{}, err
	}

	return scaffold{
		projectFolderPath: projectFolderPath,
		projectName:       projectName,
		languageScaffold:  languageScaffold,
		language:          language,
	}, nil
}

type languageScaffold interface {
	CreateProjectTemplate(projectFolderPath string, projectName string) error
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
	log.Print("Creating " + scaffold.language + " project template")
	return scaffold.languageScaffold.CreateProjectTemplate(scaffold.projectFolderPath, scaffold.projectName)
}

func (scaffold scaffold) CreateGitIgnore() error {
	log.Print("Creating: .gitignore")
	return scaffold.languageScaffold.CreateGitIgnore(scaffold.projectFolderPath)
}

func (scaffold scaffold) CreateReadme() error {
	log.Print("Creating: Readme.md")
	return util.CreateFile(scaffold.projectFolderPath, "README.md", "#"+scaffold.projectName)
}

func (scaffold scaffold) CreateEditorConfig() error {
	log.Print("Creating: .editorconfig")
	return util.CreateFile(scaffold.projectFolderPath, ".editorconfig", "")
}

func (scaffold scaffold) CreateDockerFile() error {
	log.Print("Creating: Dockerfile")
	return scaffold.languageScaffold.CreateDockerFile(scaffold.projectFolderPath, scaffold.projectName)
}

func (scaffold scaffold) CreateDockerIgnore() error {
	log.Print("Creating: .dockerignore")
	return scaffold.languageScaffold.CreateDockerIgnore(scaffold.projectFolderPath)
}

func (scaffold scaffold) CreateDeploymentFile(dockerRegistryUrl string) error {
	deploymentFilePath := filepath.Join(scaffold.projectFolderPath, "deployment_files")
	err := os.Mkdir(deploymentFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	deployYaml := template.GenericDeployYaml(scaffold.projectName, dockerRegistryUrl)

	return util.CreateFile(deploymentFilePath, "deploy.yaml", deployYaml)
}

func getLanguageScaffold(language string) (languageScaffold, error) {
	var scaffold languageScaffold

	switch language {
	case languageEnum.GoLang:
		scaffold = goLangScaffold{}
	default:
		return goLangScaffold{}, errors.New("Selected Language not supported, please see configuration options in README.md: " + language)
	}

	return scaffold, nil
}

func CreateProjectFolder(projectName string) (string, error) {
	log.Print("Creating Project folder: " + projectName)
	workingDirectory, err := os.Getwd()
	util.HandleFatalError(err)

	projectFolderPath := filepath.Join(workingDirectory, projectName)
	err = os.Mkdir(projectFolderPath, os.ModePerm)

	return projectFolderPath, err
}
