package kube

import (
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/template"
)

func DeployAwsSqs() error {
	k8, err := kubernetes.New("--context docker-for-desktop --namespace kube-system")
	if err != nil {
		return err
	}

	sqsTemplate := template.SqsDeployYaml()

	err = k8.Apply(sqsTemplate)
	if err != nil {
		return err
	}

	return nil
}