package rabbitmq

import (
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"strings"
)

func PromptRabbitSetup(envName string) error {
	user, err := util.PromptString("User Name")
	if err != nil {
		return err
	}

	password, err := util.PromptPassword("Password")
	if err != nil {
		return err
	}

	templateContent, err := generateDeployYaml(user, password)
	if err != nil {
		return err
	}

	err = template.WriteTempalte(templateContent, envName, "rabbitmq", "deploy.yaml")
	if err != nil {
		return err
	}

	return nil
}

func generateDeployYaml(user string, password string) (string, error) {
	templateContent := template.GetRabbitDeployYaml()
	templateContent = strings.Replace(templateContent, "{user}", user, -1)
	templateContent = strings.Replace(templateContent, "{password}", password, -1)

	return templateContent, nil
}