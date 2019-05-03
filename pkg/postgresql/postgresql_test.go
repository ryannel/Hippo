package postgresql

import "testing"

func TestSetupTemplate(t *testing.T) {
	_, err := generateDeployYaml("testDbName", "testUser", "testPassword")
	if err != nil {
		t.Fail()
	}
}

func TestCreateDb(t *testing.T) {
	err := CreateDb("dbName", "user", "password")
	if err != nil {
		t.Fail()
	}
}