package kube

import (
"github.com/ryannel/hippo/pkg/logger"
"github.com/ryannel/hippo/pkg/template"
)

func DeployDynamoDb() error {
	k8, err := createK8LocalInstance()
	if err != nil {
		return err
	}

	rabbitTemplate := template.DynamoDbDeployYaml()

	logger.Log("Creating dynamodb kubernetes instance")
	err = k8.Apply(rabbitTemplate)
	if err != nil {
		return err
	}

	logger.Log("dynamodb instance created on: localhost:8000")

	return nil
}
