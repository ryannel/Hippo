package setup

import (
	"github.com/ryannel/hippo/pkg/template"
	"log"
)

func SetupLocalRabbit() error {

	k8, err := createK8LocalInstance()
	if err != nil {
		return err
	}

	rabbitTemplate := template.RabbitDeployYaml("admin", "password")

	log.Print("Creating rabbit kubernetes instance")
	err = k8.Apply(rabbitTemplate)
	if err != nil {
		return err
	}

	log.Print("Creating Rabbit Secret `shared-rabbitmq`")
	secretName := "shared-rabbitmq"
	secrets := map[string]string{
		"RABBITMQ_HOST":     "rabbitmq",
		"RABBITMQ_PORT":     "5672",
		"RABBITMQ_USERNAME": "admin",
		"RABBITMQ_PASSWORD": "password",
		"RABBITMQ_VHOST":    "",
	}

	_ = k8.DeleteSecret(secretName)

	return k8.CreateSecret(secretName, secrets)

}
