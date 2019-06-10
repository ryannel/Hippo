package kubernetes

import (
	"errors"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
)

func New(context string) (Kubernetes, error) {
	command, err := getKubectlCommand()
	if err != nil {
		return Kubernetes{}, err
	}

	return Kubernetes{command + " " + context}, nil
}

type Kubernetes struct {
	command string
}

func (k8 *Kubernetes) GetPodName(appName string) (string, error) {
	command := k8.command + ` get pods --selector app=` + appName + ` --output jsonpath={.items..metadata.name}`
	logger.Command("Getting pod name: " + command)

	podName, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return string(podName), nil
}

func (k8 *Kubernetes) Apply(deployYaml string) error {
	file, err := createTmpFile(deployYaml)
	if err != nil {
		return err
	}

	command := k8.command + " apply -f " + file.Name()
	logger.Command("Applying Yaml file: " + command)

	result, err := util.ExecStringCommand(command)
	logger.Log(result)

	_ = file.Close()
	_ = os.Remove(file.Name())

	return err
}

func (k8 *Kubernetes) CreateSecret(secretName string, secrets map[string]string) error {
	command := k8.command + " create secret generic " + secretName

	for name, value := range secrets {
		command = command + " --from-literal=" + name + "=" + value
	}

	logger.Command("Creating Secrets: " + command)
	result, err := util.ExecStringCommand(command)
	if err != nil {
		return err
	}
	logger.Log(result)
	return nil
}

func (k8 *Kubernetes) DeleteSecret(secretName string) error {
	command := k8.command + " delete secret " + secretName
	logger.Command("Deleting Secret if Exists: " + command)
	result, err := util.ExecStringCommand(command)
	if err != nil {
		return err
	}
	logger.Log(result)
	return nil
}

func createTmpFile(deployYaml string) (*os.File, error) {
	file, err := ioutil.TempFile("", "hippoKubeDeploy")
	if err != nil {
		return file, err
	}

	_, err = file.WriteString(deployYaml)
	if err != nil {
		return file, err
	}

	return file, nil
}

func getKubectlCommand() (string, error) {
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
