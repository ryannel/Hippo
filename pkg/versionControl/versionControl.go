package versionControl

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ryannel/hippo/pkg/enum/versionControlProviders"
	"github.com/ryannel/hippo/pkg/util"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

func New(provider string, namespace string, project string, repository string, username string, password string) (VersionControl, error) {
	vcUrl, err := getVcUrl(provider, namespace, project, repository)
	if err != nil {
		return VersionControl{}, err
	}

	return VersionControl{provider, namespace, project, repository, username, password, vcUrl}, nil
}

type VersionControl struct {
	provider      string
	namespace     string
	project       string
	repository    string
	username      string
	password      string
	repositoryUrl string
}

func (vcs *VersionControl) Init() error {
	command := "git init"
	log.Print("Creating Git repo: " + command)
	_, err := util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) Validate() error {
	gitPath := filepath.Join(".", ".git")
	exists, err := util.PathExists(gitPath)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("not a git repository. Please run `hippo setup git`")
	}

	return nil
}

func (vcs *VersionControl) GetBranch() (string, error) {
	err := vcs.Validate()
	if err != nil {
		return "", err
	}

	command := "git rev-parse --abbrev-ref HEAD"
	log.Print("Getting branch name: " + command)
	branch, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return util.StripNewLine(branch), nil
}

func (vcs *VersionControl) GetBranchReplaceSlash() (string, error) {
	branch, err := vcs.GetBranch()
	if err != nil {
		return "", err
	}

	branch = strings.Replace(branch, "/", "_", -1)

	return branch, nil
}

func (vcs *VersionControl) GetCommit() (string, error) {
	err := vcs.Validate()
	if err != nil {
		return "", err
	}

	command := "git rev-parse HEAD"
	log.Print("Getting Commit: " + command)
	commit, err := util.ExecStringCommand(command)
	if err != nil {
		return "", err
	}

	return util.StripNewLine(commit), nil
}

func (vcs *VersionControl) CreateRepository() error {
	url := "https://dev.azure.com/fabrikam/_apis/git/repositories?api-version=5.0"
	log.Print("Creating remote repository: " + vcs.repositoryUrl)

	request := struct {
		name    string
		project struct {
			id string
		}
	}{vcs.repository, struct{ id string }{""}}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(vcs.username, vcs.password)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (vcs *VersionControl) SetOrigin() error {
	err := vcs.Validate()
	if err != nil {
		return err
	}

	command := "git remote add origin " + vcs.repositoryUrl
	log.Print("Adding Git Origin: " + command)
	_, err = util.ExecStringCommand(command)
	return err
}

func (vcs *VersionControl) CreateCommit(message string) error {
	err := vcs.Validate()
	if err != nil {
		return err
	}

	command := `git commit -m "` + message + `"`
	log.Print("Creating Git commit: " + command)
	_, err = exec.Command("git", "commit", "-m", `"`+message+`"`).Output()
	return err
}

func (vcs *VersionControl) TrackAllFiles() error {
	err := vcs.Validate()
	if err != nil {
		return err
	}

	command := "git add ."
	log.Print("Add all files to tracking: " + command)
	_, err = util.ExecStringCommand(command)
	return err
}

func getVcUrl(provider string, namespace string, project string, repository string) (string, error) {
	var vcUrl string

	switch provider {
	case versionControlProviders.Azure:
		vcUrl = buildAzureVcUrl(namespace, project, repository)
	default:
		return "", errors.New("Invalid Version control Provider")
	}

	return vcUrl, nil
}

func buildAzureVcUrl(vcNamespace string, vcProject string, vcRepository string) string {
	return "http://" + vcNamespace + ".visualstudio.com/" + vcProject + "/_git/" + vcRepository
}
