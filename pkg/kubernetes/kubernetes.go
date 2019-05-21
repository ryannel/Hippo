package kubernetes

import (
	"errors"
	"github.com/ryannel/hippo/pkg/util"
	"io/ioutil"
	"log"
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

func (k8 *Kubernetes) Deploy(deployYamlPath string) error {
	command := k8.command + " apply --record=false -f " + deployYamlPath
	log.Print("Deploying pod: " + command)

	result, err := util.ExecStringCommand(command)
	if err != nil {
		return err
	}

	print(result)

	return nil
}

func (k8 *Kubernetes) GetPodName(appName string) (string, error) {
	command := k8.command + ` get pods --selector app=` + appName + ` --output jsonpath={.items..metadata.name}`
	log.Print("Getting pod name: " + command)

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
	log.Print("Applying Yaml file: " + command)

	result, err := util.ExecStringCommand(command)
	log.Print(result)

	_ = file.Close()
	_ = os.Remove(file.Name())

	return err
}

func (k8 *Kubernetes) CreateSecret(secretName string , secrets map[string]string) error {
	command := k8.command + " create secret generic " + secretName

	for name, value := range secrets {
		command = command + " --from-literal=" + name + "="+`"` + value + `"`
	}

	log.Print("Creating Secrets: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func (k8 *Kubernetes) DeleteSecret(secretName string) error {
	command := k8.command + "delete secret " + secretName
	log.Print("Deleting Secret if Exists: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func createTmpFile(deployYaml string) (*os.File, error){
	file, err := ioutil.TempFile("", "psqlKubeDeploy")
	if err != nil {
		return file, err
	}

	_, err = file.WriteString(deployYaml)
	if err != nil {
		return file, err
	}

	return file, nil
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

func GetContext(contextName string, contexts map[string]string) (string, error){
	var err error = nil

	kubeContext := contexts[contextName]
	if len(kubeContext) == 0 {
		err = errors.New("env name (" + contextName + ") not found in hippo.yaml KubernetesContext list" )
	}

	return kubeContext, err
}


