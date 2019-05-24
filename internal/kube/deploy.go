package kube

import (
	"errors"
	"github.com/ryannel/hippo/pkg/configManager"
	"github.com/ryannel/hippo/pkg/kubernetes"
	"github.com/ryannel/hippo/pkg/util"
	"github.com/ryannel/hippo/pkg/versionControl"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Deploy(envName string) error {
	projectFolder, err := os.Getwd()
	if err != nil {
		return err
	}

	deployYamlPath := filepath.Join(projectFolder, "deployment_files", "deploy.yaml")

	exists, err := util.PathExists(deployYamlPath)
	if !exists || err != nil {
		return errors.New("deployment files do not exist. run `hippo setup kubernetes` to create them: " + deployYamlPath)
	}

	confManager, err := configManager.New("hippo.yaml")
	if err != nil {
		return err
	}

	config := confManager.GetConfig()

	kubeEnv := config.KubernetesContexts[envName]
	if len(kubeEnv) == 0 {
		return errors.New("not a valid kubernetes context. Please ensure the context name exists in hippo.yaml. Run `hippo setup kubernetes` to configure")
	}

	k8, err := kubernetes.New(kubeEnv)
	if err != nil {
		return err
	}

	vcs, err := versionControl.New()
	if err != nil {
		return errors.New("unable to find git. Please run `git init` and create a commit")
	}

	commitTag, err := vcs.GetCommit()
	if err != nil {
		return errors.New("unable to find latest commit. Please ensure that this branch contains at least one commit")
	}

	template, err := ioutil.ReadFile(deployYamlPath)
	if err != nil {
		return err
	}
	deployYaml := string(template)

	log.Print("Setting deploy.yaml {commit} to: " + commitTag)
	deployYaml = strings.Replace(deployYaml, "{commit}", commitTag, -1)
	deployYaml = strings.Replace(deployYaml, "{timestamp}", time.Now().Format(time.RFC3339), -1)

	return k8.Apply(deployYaml)
}