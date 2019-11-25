package kubernetes

import (
	"errors"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func (k8 *Kubernetes) GetSecrets(namespace string) ([]string, error) {
	command := k8.command + " -n " +namespace+ " get secrets"
	logger.Command("Fetching secrets: " + command)
	result, err := util.ExecStringCommand(command)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")[1:]

	var secretNames []string
	for _, row := range rows {
		name := strings.Split(row, " ")[0]
		if name != "" {
			secretNames = append(secretNames, name)
		}
	}

	return secretNames, err
}

func (k8 *Kubernetes) GetEksAdminToken() (string, error){
	secretNames, err := k8.GetSecrets("kube-system")
	if err != nil {
		return "", err
	}

	var eksAdminSecretName string
	for _, name := range secretNames {
		matched, err := filepath.Match("eks-admin-token-*", name)
		if err != nil {
			return "", err
		}
		if matched {
			eksAdminSecretName = name
		}
	}

	command := k8.command + " -n kube-system describe secret " + eksAdminSecretName
	logger.Command("Getting EKS Admin Token: " + command)
	result, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	rows := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	var token string
	for _, row := range rows {
		matched, err := filepath.Match("token:*", row)
		if err != nil {
			return "", err
		}
		if matched {
			token = strings.TrimSpace(row[6:])
			break
		}
	}

	if token == "" {
		return "", errors.New("unable to find token")
	}

	return token, nil
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
