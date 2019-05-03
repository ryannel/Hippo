package rabbitmq

import "testing"

func TestSetupTemplate(t *testing.T) {
	_, err := generateDeployYaml( "testUser", "testPassword")

	if err != nil {
		t.Fail()
	}
}