package elasticSearch

import (
	"errors"
	awsElasticSearch "github.com/aws/aws-sdk-go/service/elasticsearchservice"
)

type ElasticSearch struct {
	Sdk *awsElasticSearch.ElasticsearchService
}

func (elastic *ElasticSearch) GetDomains() ([]string, error) {
	domains, err := elastic.Sdk.ListDomainNames(nil)
	if err != nil {
		return nil, err
	}
	var domainNames []string
	for _, domain := range domains.DomainNames {
		domainNames = append(domainNames, *domain.DomainName)
	}
	return domainNames, nil
}

func (elastic *ElasticSearch) GetVpcEndpoint(domainName string) (string, error) {
	input := awsElasticSearch.DescribeElasticsearchDomainInput{DomainName: &domainName}

	domain, err := elastic.Sdk.DescribeElasticsearchDomain(&input)
	if err != nil {
		return "", err
	}

	var vpcDomain string
	for key, endpoint := range domain.DomainStatus.Endpoints {
		if key == "vpc" {
			vpcDomain = *endpoint
			break
		}
	}

	if vpcDomain == "" {
		return "", errors.New("no elastic vpc endpoint found")
	}

	return vpcDomain, nil
}