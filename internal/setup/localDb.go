package setup

import (
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/postgresql"
	"github.com/ryannel/hippo/pkg/template"
	"log"
)

func SetupLocalDb() error {
	config, err := configuration.New()
	if err != nil {
		return err
	}

	k8, err := createK8LocalInstance()
	if err != nil {
		return err
	}

	err = createPostgresContainer(k8)
	if err != nil {
		return err
	}

	log.Print()
	psql, err := connectToPsql()
	if err != nil {
		return err
	}

	err = createDevDb(psql, config.ProjectName)
	if err != nil {
		log.Print(err)
	}

	err = createDbUser(psql, config.ProjectName)
	if err != nil {
		log.Print(err)
	}

	log.Print()
	err = setDevDbSecret(k8, config.ProjectName)
	if err != nil {
		return err
	}

	return nil
}

func createK8LocalInstance()(kubernetes.Kubernetes, error) {
	k8, err := kubernetes.New("--context docker-for-desktop --namespace default")
	if err != nil {
		return kubernetes.Kubernetes{}, err
	}
	return k8, nil
}

func createPostgresContainer(k8 kubernetes.Kubernetes) error {
	psqlTemplate := template.PostgresDeployYaml("postgres", "postgres", "postgres")

	log.Print("Creating Postgresql kubernetes instance")
	err := k8.Apply(psqlTemplate)
	if err != nil {
		return err
	}

	log.Print("Creating Root DB Secret `shared-postgres`")
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
	log.Print("Connecting to DB instance")
	psql, err := postgresql.New("localhost", 5432, "postgres", "postgres", "postgres")
	if err != nil {
		return postgresql.Postgresql{}, err
	}
	return psql, nil
}

func createDevDb(psql postgresql.Postgresql, projectName string) error {
	log.Print("Creating dev db: `" + projectName + "` with owner `" + projectName + "`")
	return psql.CreateDb(projectName, projectName)
}

func createDbUser(psql postgresql.Postgresql, projectName string) error {
	log.Print("Creating dev user: `" + projectName + "` with password `" + projectName + "`")
	return psql.CreateUser(projectName, projectName)
}

func setDevDbSecret(k8 kubernetes.Kubernetes, projectName string) error {
	log.Print("Creating Dev DB Secret `" + projectName + "`")
	secretName := projectName
	secrets := map[string]string{
		"POSTGRES_HOST":     projectName,
		"POSTGRES_DB":       projectName,
		"POSTGRES_USER":     projectName,
		"POSTGRES_PASSWORD": projectName,
	}

	_ = k8.DeleteSecret(secretName)
	return k8.CreateSecret(secretName, secrets)
}
