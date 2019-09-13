package aws

import (
	"errors"
	"github.com/ryannel/hippo/pkg/aws"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"time"
)

func ConnectElasticSearch(region string, profile string) error {
	connection, err := aws.Connect(region)
	if err != nil {
		return err
	}

	result, err := connection.Login(profile)
	if err != nil {
		return err
	}
	logger.Info(result)

	workerId, err := connection.EC2.GetRunningWorkerInstanceId()
	if err != nil {
		return err
	}

	domains, err := connection.ElasticSearch.GetDomains()
	if err != nil {
		return err
	}

	endpoint, err := connection.ElasticSearch.GetVpcEndpoint(domains[0])
	if err != nil {
		return err
	}

	pemFilePath, err := getCertificatePath()
	if err != nil {
		return err
	}

	command := "ssh -i " + pemFilePath + " ec2-user@" + workerId + " -N -L 30443:"+endpoint+":443"
	logger.Command("Starting SSH tunnel: " + command)
	go func() {
		_, err = util.ExecStringCommand(command)
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	time.Sleep(5 * time.Second)
	logger.Info("SSH Tunnel Started")
	util.Openbrowser("https://localhost:30443/_plugin/kibana")

	util.WaitForever()
	return nil
}

func getCertificatePath() (string, error){
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