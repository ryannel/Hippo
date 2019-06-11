package kube

import (
	"fmt"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/template"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"os/exec"
	"runtime"
)

func Ui() error {
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
	openbrowser("http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/")

	logger.Command("starting proxy server: kubectl proxy")
	result, err := util.ExecStringCommand("kubectl proxy")
	if err != nil {
		return err
	}
	logger.Log(result)

	return nil
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
