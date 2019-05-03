package postgresql

import (
	"hippo/pkg/kubernetes"
	"hippo/pkg/template"
	"hippo/pkg/util"
	"os/exec"
	"strings"
)

func PromptPostgresSetup(envName string) error {
	dbName, err := util.PromptString("Database Name")
	if err != nil {
		return err
	}

	user, err := util.PromptString("User Name")
	if err != nil {
		return err
	}

	password, err := util.PromptPassword("Password")
	if err != nil {
		return err
	}

	templateContent, err := generateDeployYaml(dbName, user, password)
	if err != nil {
		return err
	}

	err = template.WriteTempalte(templateContent, envName, "postgresql", "deploy.yaml")
	if err != nil {
		return err
	}

	return nil
}

func generateDeployYaml(dbName string, user string, password string) (string, error) {
	templateContent := template.GetPostgresDeployYaml()
	templateContent = strings.Replace(templateContent, "{dbName}", dbName, -1)
	templateContent = strings.Replace(templateContent, "{user}", user, -1)
	templateContent = strings.Replace(templateContent, "{password}", password, -1)

	return templateContent, nil
}

func CreateDb(dbName string, user string, password string) error {
	podName, err := kubernetes.GetPodName("postgresql")
	if err != nil {
		return err
	}

	println(`Creating development db: kubectl.exe exec -it ` + podName + ` -- bash -c "echo \"CREATE USER \\\"` + user + `\\\" WITH PASSWORD '` + password + `'; CREATE DATABASE \\\"` + dbName + `\\\" WITH OWNER \\\"` + user + `\\\" ENCODING utf8\" | psql -U postgres -f -"`)

	_, err = exec.Command("kubectl.exe",  "exec", "-it", podName, "--", "bash", "-c", `"echo hello"`).Output()
	//_, err = exec.Command("kubectl.exe",  "exec", "-it", podName, "--", "bash", "-c", `"echo \"CREATE USER \\\"` + user + `\\\" WITH PASSWORD '` + password + `'; CREATE DATABASE \\\"` + dbName + `\\\" WITH OWNER \\\"` + user + `\\\" ENCODING utf8\" | psql -U postgres -f -"`).Output()
	if err != nil {
		return err
	}

	return nil
}