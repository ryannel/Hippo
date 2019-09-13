package kube

import (
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
)

func DeployDashboard() error {
	k8, err := kubernetes.New("--context docker-for-desktop --namespace kube-system")
	if err != nil {
		return err
	}

	dashboardTemplate := template.KubeDashboard()

	err = k8.Apply(dashboardTemplate)
	if err != nil {
		return err
	}

	logger.Command("Opening Browser")
	util.Openbrowser("http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/")

	logger.Command("starting proxy server: kubectl proxy")
	result, err := util.ExecStringCommand("kubectl proxy")
	if err != nil {
		return err
	}
	logger.Log(result)

	return nil
}
