package azure

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Provider struct {
	Namespace  string
	Project    string
	Repository string
	Username   string
	Password   string
}

func (azure *Provider) CreateRepository() error {
	url := "https://dev.azure.com/" + azure.Namespace + "/" + azure.Project + "/_apis/git/repositories?api-version=5.0"
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
		Name: azure.Repository,
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
	req.SetBasicAuth(azure.Username, azure.Password)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 409: err = errors.New("409: repository already exist")
	case 201: err = errors.New(response.Status)
	default: err = nil
	}

	return err
}

//func (azure *Provider) CreateBuildPipeline() error {
//	type requestObject struct {
//		authoredBy struct {
//
//		}
//	}
//
//	url := "POST https://dev.azure.com/"+azure.Namespace+"/"+azure.Project+"/_apis/build/definitions?api-version=5.0"
//	return nil
//}

func (azure *Provider) getProjectId() (string, error) {
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

	url := "https://dev.azure.com/" + azure.Namespace + "/_apis/projects/" + azure.Project + "?api-version=5.0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(azure.Username, azure.Password)

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

func (azure *Provider) GetRepositoryUrl() string {
	return "http://" + azure.Namespace + ".visualstudio.com/" + azure.Project + "/_git/" + azure.Repository
}
