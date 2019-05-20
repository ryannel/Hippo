package postgresql

import (
	"database/sql"
)

//
//import (
//	_ "github.com/lib/pq"
//	"github.com/ryannel/hippo/pkg/template"
//	"github.com/ryannel/hippo/pkg/util"
//	"strings"
//)
//
//func PromptPostgresSetup(envName string) error {
//	dbName, err := util.PromptString("Database Name")
//	if err != nil {
//		return err
//	}
//
//	user, err := util.PromptString("User Name")
//	if err != nil {
//		return err
//	}
//
//	password, err := util.PromptPassword("Password")
//	if err != nil {
//		return err
//	}
//
//	templateContent, err := generateDeployYaml(dbName, user, password)
//	if err != nil {
//		return err
//	}
//
//	err = template.WriteTempalte(templateContent, envName, "postgresql", "deploy.yaml")
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func generateDeployYaml(dbName string, user string, password string) (string, error) {
//	templateContent := template.GetPostgresDeployYaml()
//	templateContent = strings.Replace(templateContent, "{dbName}", dbName, -1)
//	templateContent = strings.Replace(templateContent, "{user}", user, -1)
//	templateContent = strings.Replace(templateContent, "{password}", password, -1)
//
//	return templateContent, nil
//}
//

func New(host string, dbName string, user string, password string) (Postgresql, error) {
	connStr := "User ID=" + user + ";Password=" + password + ";Host=" + host + ";Port=5432;Database=" + dbName + ";"
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return Postgresql{}, err
	}

	return Postgresql{connStr, connection}, nil
}

type Postgresql struct {
	connectionString string
	connection *sql.DB
}

func (psql *Postgresql) CreateUser(username string, password string) error {
	_, err := psql.connection.Exec(`CREATE USER "` + username + ` WITH PASSWORD '`+ password + `';`)
	return err
}

func (psql *Postgresql) CreateDb(dbName string, owner string) error {
	_, err := psql.connection.Exec(`CREATE DATABASE "` + dbName + ` WITH OWNER "`+ owner + `" ENCODING utf8;`)
	return err
}

func (psql *Postgresql) Close() error {
	return psql.connection.Close()
}