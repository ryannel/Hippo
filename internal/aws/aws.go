package aws

import (
	"errors"
	"github.com/ryannel/hippo/pkg/aws"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func ConnectElasticSearch(region string, profile string) error {
	err := os.Setenv("AWS_PROFILE", profile)
	if err != nil {
		return err
	}

	result, err := aws.Login(profile)
	if err != nil {
		return err
	}
	logger.Info(result)

	connection, err := aws.New(region)
	if err != nil {
		return err
	}

	logger.Info("Finding running worker instance")
	workerId, err := connection.EC2.GetRunningWorkerInstanceId()
	if err != nil {
		return err
	}
	logger.Info("Using worker instance: " + workerId)

	logger.Info("Finding Elastic Search domains")
	domains, err := connection.ElasticSearch.GetDomains()
	if err != nil {
		return err
	}

	logger.Info("Finding VPC end point for domain: " + domains[0])
	endpoint, err := connection.ElasticSearch.GetVpcEndpoint(domains[0])
	if err != nil {
		return err
	}
	logger.Info("Using VPC endpoint: " + endpoint)

	pemFilePath, err := getCertificatePath()
	if err != nil {
		return err
	}

	command := "ssh -i " + pemFilePath + " ec2-user@" + workerId + " -N -L 30443:" + endpoint + ":443"
	logger.Command("Executing SSH tunnel: " + command)

	cmd, errCh := execAsyncCommand(command)
	select {
	case err, errored := <-errCh:
		if errored {
			return err
		}
	default:
	}

	logger.Info("SSH starting...")
	time.Sleep(5 * time.Second)
	logger.Info("SSH tunnel created")
	util.Openbrowser("https://localhost:30443/_plugin/kibana")

	return cmdAwaitInterrupt(cmd, errCh, "Shutting down SSH tunnel")
}

func ConnectPostgres(region string, profile string) error {
	err := os.Setenv("AWS_PROFILE", profile)
	if err != nil {
		return err
	}

	result, err := aws.Login(profile)
	if err != nil {
		return err
	}
	logger.Info(result)

	connection, err := aws.New(region)
	if err != nil {
		return err
	}

	logger.Info("Finding running worker instance")
	workerId, err := connection.EC2.GetRunningWorkerInstanceId()
	if err != nil {
		return err
	}
	logger.Info("Using worker instance: " + workerId)

	logger.Info("Finding RDS instances")
	instances, err := connection.RDS.GetInstances()
	if err != nil {
		return err
	}

	logger.Info("Finding endpoint for RDS instance: " + instances[0])
	endpoint, err := connection.RDS.GetEndpoint(instances[0])

	pemFilePath, err := getCertificatePath()
	if err != nil {
		return err
	}

	command := "ssh -i " + pemFilePath + " ec2-user@" + workerId + " -N -L 35432:" + endpoint.Address + ":" + strconv.Itoa(endpoint.Port) + ""
	logger.Command("Executing SSH tunnel: " + command)

	cmd, errCh := execAsyncCommand(command)
	select {
	case err, errored := <-errCh:
		if errored {
			return err
		}
	default:
	}

	logger.Info("SSH starting...")
	time.Sleep(3 * time.Second)
	logger.Info("SSH tunnel created. Db Exposed on: `localhost:35432`")

	return cmdAwaitInterrupt(cmd, errCh, "Shutting down SSH tunnel")
}

func SetContext(contextName string) error {
	if contextName != "local" {
		result, err := aws.Login(contextName)
		if err != nil {
			return err
		}
		logger.Info(result)
	}

	command := "kubectl config use-context " + contextName
	logger.Command("using kubectl context `" + contextName + "`: " + command)
	_, err := util.ExecStringCommand(command)
	if err != nil {
		return err
	}

	logger.Log("AWS context switched to: " + contextName)
	return nil
}

func ConnectDashboard(profile string) error {
	result, err := aws.Login(profile)
	if err != nil {
		return err
	}
	logger.Info(result)

	k8, err := kubernetes.New("--context="+profile)
	if err != nil {
		return err
	}
	eksAdminToken, err := k8.GetEksAdminToken()

	command := "kubectl proxy --port=8010 --context=" + profile
	logger.Command("Proxying the dashboard port: " + command)
	cmd, errCh := execAsyncCommand(command)
	select {
	case err, errored := <-errCh:
		if errored {
			return err
		}
	default:
	}

	time.Sleep(3 * time.Second)
	util.Openbrowser("http://localhost:8010/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/")
	logger.Log("Use Token to login: " + eksAdminToken)

	return cmdAwaitInterrupt(cmd, errCh, "Shutting down proxy")
}

func execAsyncCommand(command string) (*exec.Cmd, chan error) {
	errCh := make(chan error)
	cmdCh := make(chan *exec.Cmd)

	go func() {
		args := strings.Fields(command)
		cmd := exec.Command(args[0], args[1:]...)
		cmdCh <- cmd

		result, err := cmd.Output()
		logger.Info(string(result))
		errCh <- err
	}()

	return <-cmdCh, errCh
}

func cmdAwaitInterrupt(cmd *exec.Cmd, errCh chan error, shutdownMessage string) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	var err error
	select {
	case <-errCh:
		logger.Error((<-errCh).Error())
		err = cmd.Process.Kill()
	case <-sigCh:
		logger.Info(shutdownMessage)
		err = cmd.Process.Kill()
		if err != nil {
			logger.Error(err.Error())
		}

	}

	return err
}

func getCertificatePath() (string, error) {
	config, err := configuration.New()
	if err != nil {
		return "", err
	}

	if config.ConfigPath == "" {
		return "", errors.New("no hippo.yaml found in path. Please run `hippo configure`")
	}

	if config.Aws.CertificatePath == "" {
		certificatePath, err := util.PromptString("AWS Pem file path")
		if err != nil {
			return "", err
		}
		config.Aws.CertificatePath = certificatePath
		err = config.SaveConfig()
		if err != nil {
			return "", err
		}
	}

	return config.Aws.CertificatePath, nil
}
