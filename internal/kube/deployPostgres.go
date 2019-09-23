package kube

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/postgresql"
	"github.com/ryannel/hippo/pkg/template"
)

func DeployPostgres() error {
	config, err := configuration.New()
	if err != nil {
		return err
	}

	if config.ConfigPath == "" {
		return errors.New("no hippo.yaml found in path. Please run `hippo configure`")
	}

	k8, err := createK8LocalInstance()
	if err != nil {
		return err
	}

	err = createPostgresContainer(k8)
	if err != nil {
		return err
	}

	psql, err := connectToPsql()
	if err != nil {
		return err
	}

	err = createDbUser(psql, config.ProjectName)
	if err != nil {
		logger.Warn(err.Error())
	}

	err = createDevDb(psql, config.ProjectName)
	if err != nil {
		logger.Warn(err.Error())
	}

	err = setDevDbSecret(k8, config.ProjectName)
	if err != nil {
		return err
	}

	return nil
}

func createK8LocalInstance() (kubernetes.Kubernetes, error) {
	k8, err := kubernetes.New("")
	if err != nil {
		return kubernetes.Kubernetes{}, err
	}
	return k8, nil
}

func createPostgresContainer(k8 kubernetes.Kubernetes) error {
	psqlTemplate := template.PostgresDeployYaml("postgres", "postgres", "postgres")

	logger.Command("Creating Postgresql kubernetes instance")
	err := k8.Apply(psqlTemplate)
	if err != nil {
		return err
	}

	logger.Command("Creating Root DB Secret `shared-postgres`")
	secretName := "shared-postgres"
	secrets := map[string]string{
		"POSTGRES_HOST":     "postgres",
		"POSTGRES_DB":       "postgres",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
	}

	_ = k8.DeleteSecret(secretName)
	return k8.CreateSecret(secretName, secrets)
}

func connectToPsql() (postgresql.Postgresql, error) {
	logger.Command("Connecting to DB instance")
	psql, err := postgresql.New("localhost", 5432, "postgres", "postgres", "postgres")
	if err != nil {
		return postgresql.Postgresql{}, err
	}
	return psql, nil
}

func createDevDb(psql postgresql.Postgresql, projectName string) error {
	logger.Command("Creating dev db: `" + projectName + "` with owner `" + projectName + "`")
	return psql.CreateDb(projectName, projectName)
}

func createDbUser(psql postgresql.Postgresql, projectName string) error {
	logger.Command("Creating dev user: `" + projectName + "` with password `" + projectName + "`")
	return psql.CreateUser(projectName, projectName)
}

func setDevDbSecret(k8 kubernetes.Kubernetes, projectName string) error {
	logger.Command("Creating Dev DB Secret `" + projectName + "`")
	secretName := projectName + "-db"
	secrets := map[string]string{
		"POSTGRES_HOST":     "postgresql",
		"POSTGRES_DBNAME":   projectName,
		"POSTGRES_USERNAME": projectName,
		"POSTGRES_PASSWORD": projectName,
	}

	_ = k8.DeleteSecret(secretName)
	return k8.CreateSecret(secretName, secrets)
}
