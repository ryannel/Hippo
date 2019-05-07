package kubernetes

import (
	"errors"
	"hippo/pkg/environment"
	"hippo/pkg/util"
	"log"
	"os/exec"
)

func New(envName string, config environment.EnvConfig) (Kubernetes, error) {
	command, err := getKubeContext(envName, config)
	if err != nil {
		return Kubernetes{}, err
	}

	return Kubernetes{command}, nil
}

type Kubernetes struct {
	contextCommand string
}

func (k8 *Kubernetes) Deploy(deployYamlPath string) error {
	command := k8.contextCommand + "apply --record=false -f " + deployYamlPath
	log.Print("Deploying pod: " + command)

	result, err := util.ExecStringCommand(command)
	if err != nil {
		return err
	}

	print(result)

	return nil
}

func (k8 *Kubernetes) GetPodName(appName string) (string, error) {
	command := k8.contextCommand + ` get pods --selector app=` + appName + ` --output jsonpath={.items..metadata.name}`
	log.Print("Getting pod name: " + command)

	podName, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return string(podName), nil
}

func getKubectlCommand() (string, error){
	_, err := exec.LookPath("kubectl")
	if err == nil {
		return "kubectl", nil
	}

	_, err = exec.LookPath("kubectl.exe")
	if err == nil {
		return "kubectl.exe", nil
	}

	return "", errors.New("kubectl command not found")
}

func getContextArgs(envName string, config environment.EnvConfig) (string, error) {
	context := config.Environments[envName]
	if len(context) == 0 {
		return context, errors.New("env name (" + envName + ") not found in hippo.yaml Environments list" )
	}
	return context, nil
}

func getKubeContext(envName string, config environment.EnvConfig) (string, error){
	kubectl, err := getKubectlCommand()
	if err != nil {
	    return "", err
	}

	contextArgs, err := getContextArgs(envName, config)
	if err != nil {
	    return "", err
	}

	return kubectl + " " + contextArgs + " ", nil
}


