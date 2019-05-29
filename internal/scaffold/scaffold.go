package scaffold

import (
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/scaffoldManager"
	"log"
)

func ScaffoldProject(projectName string, language string) error {
	projectFolderPath, err := scaffoldManager.CreateProjectFolder(projectName)
	if err != nil {
		return err
	}

	scaffold, err := scaffoldManager.New(projectName, projectFolderPath, language)
	if err != nil {
		return err
	}

	err = scaffold.CreateProjectTemplate()
	if err != nil {
		return err
	}

	err = scaffold.CreateEditorConfig()
	if err != nil {
		return err
	}

	err = scaffold.CreateReadme()
	if err != nil {
		return err
	}

	config, err := configuration.Create(projectFolderPath)
	if err != nil {
		return err
	}

	config.ProjectName = projectName
	config.Language = language

	err = config.SaveConfig()
	if err != nil {
		return err
	}

	log.Print("Project has been created at `./" + projectName + "`")
	log.Print("Please initialise your version control in the project folder")
	return nil
}
