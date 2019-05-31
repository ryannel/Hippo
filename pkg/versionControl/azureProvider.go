package versionControl

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type azureProvider struct {
	namespace  string
	project    string
	repository string
	username   string
	password   string
}

func (azure *azureProvider) createRepository() error {
	url := "https://dev.azure.com/" + azure.namespace + "/" + azure.project + "/_apis/git/repositories?api-version=5.0"
	log.Print("Creating remote repository: ")

	type projectObject struct {
		Id string `yaml:"Id"`
	}

	type requestObject struct {
		Name    string        `json:"Name"`
		Project projectObject `json:"Project"`
	}

	projectId, err := azure.getProjectId()
	if err != nil {
		return err
	}

	request := requestObject{
		Name: azure.repository,
		Project: projectObject{
			projectId,
		},
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(azure.username, azure.password)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return errors.New(response.Status)
	}

	return nil
}

func (azure *azureProvider) getProjectId() (string, error) {
	type projectInfoObject struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		State       string `json:"state"`
		Revision    int    `json:"revision"`
		Links       struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Collection struct {
				Href string `json:"href"`
			} `json:"collection"`
			Web struct {
				Href string `json:"href"`
			} `json:"web"`
		} `json:"_links"`
		Visibility  string `json:"visibility"`
		DefaultTeam struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"defaultTeam"`
		LastUpdateTime string `json:"lastUpdateTime"`
	}

	url := "https://dev.azure.com/" + azure.namespace + "/_apis/projects/" + azure.project + "?api-version=5.0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(azure.username, azure.password)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	decoder := json.NewDecoder(response.Body)
	var projectInfo projectInfoObject
	err = decoder.Decode(&projectInfo)
	if err != nil {
		return "", err
	}

	return projectInfo.Id, nil
}

func (azure *azureProvider) getRepositoryUrl() string {
	return "http://" + azure.namespace + ".visualstudio.com/" + azure.project + "/_git/" + azure.repository
}
