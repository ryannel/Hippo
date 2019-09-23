package kube

import (
	"errors"
	"github.com/ryannel/hippo/internal/docker"
	"github.com/ryannel/hippo/pkg/configuration"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/versionControl"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Deploy(envName string) error {
	logger.Log("Starting `" + envName + "` deployment")

	err := docker.Build()
	if err != nil {
		return err
	}

	projectFolder, err := os.Getwd()
	if err != nil {
		return err
	}

	deployYamlPath := filepath.Join(projectFolder, "deployment_files")

	exists, err := util.PathExists(deployYamlPath)
	if !exists || err != nil {
		return errors.New("deployment files do not exist. run `hippo setup kubernetes` to create them: " + deployYamlPath)
	}

	k8, err := kubernetes.New("")
	if err != nil {
		return err
	}

	templates, err := getLocalDeploymentConfigs(projectFolder, envName)
	if err != nil {
		return err
	}

	for _, template := range templates {
		err := k8.Apply(template)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLocalDeploymentConfigs(projectFolder string, environment string) ([]string, error) {
	config, err := configuration.New()
	if err != nil {
		return nil, err
	}

	if config.ConfigPath == "" {
		return nil, errors.New("no hippo.yaml found in path. Please run `hippo configure`")
	}

	var configurations []string
	environmentDeploymentFiles := filepath.Join(projectFolder, "deployment_files", environment)
	environmentConfigs, err := getTemplatesFromFolder(environmentDeploymentFiles, config)
	if err != nil {
		return nil, err
	}
	configurations = append(configurations, environmentConfigs...)

	deploymentFiles := filepath.Join(projectFolder, "deployment_files")
	deploymentConfigs, err := getTemplatesFromFolder(deploymentFiles, config)
	if err != nil {
		return nil, err
	}
	configurations = append(configurations, deploymentConfigs...)

	return configurations, nil
}

func getTemplatesFromFolder(path string, config configuration.Configuration) ([]string, error) {
	exists, err := util.PathExists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	commit, err := getCommitTag(config)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, f := range files {
		if f.IsDir() == false {
			fPath := filepath.Join(path, f.Name())
			template, err := ioutil.ReadFile(fPath)
			if err != nil {
				return nil, err
			}

			logger.Info("Setting ${COMMIT} to: " + commit + " for file: " + fPath)
			templateString := string(template)
			templateString = strings.Replace(templateString, "${COMMIT}", commit, -1)
			templateString = strings.Replace(templateString, "${TIMESTAMP}", time.Now().Format(time.RFC3339), -1)

			templates = append(templates, templateString)
		}
	}

	return templates, nil
}

func getCommitTag(config configuration.Configuration) (string, error) {
	vcs, err := versionControl.New(config.VersionControl.Provider, config.VersionControl.NameSpace, config.VersionControl.Project, config.VersionControl.Repository, config.VersionControl.Username, config.VersionControl.Password)
	if err != nil {
		return "", errors.New("unable to find git. Please run `git init` and create a commit")
	}

	commitTag, err := vcs.GetCommit()
	if err != nil {
		return "", errors.New("unable to find latest commit. Please ensure that this branch contains at least one commit")
	}

	return commitTag, nil
}