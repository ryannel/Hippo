package kube

import (
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/template"
)

func DeployRabbit() error {
	k8, err := createK8LocalInstance()
	if err != nil {
		return err
	}

	rabbitTemplate := template.RabbitDeployYaml("admin", "password")

	logger.Log("Creating rabbit kubernetes instance")
	err = k8.Apply(rabbitTemplate)
	if err != nil {
		return err
	}

	logger.Log("Creating Rabbit Secret `shared-rabbitmq`")
	secretName := "shared-rabbitmq"
	secrets := map[string]string{
		"RABBITMQ_HOST":     "rabbitmq",
		"RABBITMQ_PORT":     "5672",
		"RABBITMQ_USERNAME": "admin",
		"RABBITMQ_PASSWORD": "password",
		"RABBITMQ_VHOST":    "",
	}

	_ = k8.DeleteSecret(secretName)

	err = k8.CreateSecret(secretName, secrets)
	if err != nil {
		return err
	}

	logger.Log("Rabbit instance created on: localhost:5672")
	logger.Log("Management port: localhost:15672 User:Admin Password:password")

	return nil
}
