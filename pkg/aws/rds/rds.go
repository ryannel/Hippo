package rds

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	awsRds "github.com/aws/aws-sdk-go/service/rds"
)

type RDS struct {
	Sdk *awsRds.RDS
}

func (rds *RDS) GetInstances() ([]string, error) {
	input := awsRds.DescribeDBInstancesInput{}
	instanceDescriptions, err := rds.Sdk.DescribeDBInstances(&input)
	if err != nil {
		return nil, err
	}

	var instanceNames []string
	for _, instance := range instanceDescriptions.DBInstances {
		instanceNames = append(instanceNames, *instance.DBInstanceIdentifier)
	}

	return instanceNames, nil
}

type Endpoint struct {
	Address string
	Port int
}

func (rds *RDS) GetEndpoint(instanceName string) (Endpoint, error) {
	input := awsRds.DescribeDBInstancesInput{DBInstanceIdentifier: aws.String(instanceName)}
	instanceDescriptions, err := rds.Sdk.DescribeDBInstances(&input)
	if err != nil {
		return Endpoint{}, err
	}

	if len(instanceDescriptions.DBInstances) == 0 {
		return Endpoint{}, errors.New("no rds instance found")
	}
	if len(instanceDescriptions.DBInstances) > 1 {
		return Endpoint{}, errors.New("more than 1 rds instance selected")
	}

	endpoint := Endpoint{
		Address: *instanceDescriptions.DBInstances[0].Endpoint.Address,
		Port: int(*instanceDescriptions.DBInstances[0].Endpoint.Port),
	}

	return endpoint, nil
}