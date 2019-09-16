package aws

import (
	"errors"
	"github.com/ryannel/hippo/pkg/aws"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func ConnectElasticSearch(region string, profile string) error {
	connection, err := aws.New(region)
	if err != nil {
		return err
	}

	result, err := connection.Login(profile)
	if err != nil {
		return err
	}
	logger.Info(result)

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

	logger.Info("SSH starting...")
	time.Sleep(5 * time.Second)
	logger.Info("SSH tunnel created")
	util.Openbrowser("https://localhost:30443/_plugin/kibana")

	return cmdAwaitInterrupt(cmd, errCh, "Shutting down SSH tunnel")
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
	case <-sigCh:
		logger.Info(shutdownMessage)
		err = cmd.Process.Kill()
		if err != nil {
			logger.Error(err.Error())
		}
	case <-errCh:
		logger.Error((<-errCh).Error())
		err = cmd.Process.Kill()
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
