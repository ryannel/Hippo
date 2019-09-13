package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	awsEc2 "github.com/aws/aws-sdk-go/service/ec2"
	"path/filepath"
)

type EC2 struct {
	Sdk *awsEc2.EC2
}

func (inst *EC2) GetRunningInstances() ([]*awsEc2.Instance, error){
	filters := []*awsEc2.Filter{
		{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("running")},
		},
	}

	input := awsEc2.DescribeInstancesInput{Filters: filters}
	response, err := inst.Sdk.DescribeInstances(&input)
	if err != nil {
		return nil, err
	}

	instances, err := inst.groupInstances(response.Reservations)
	if err != nil {
		return instances, err
	}

	return instances, nil
}

func (inst *EC2) groupInstances(reservations []*awsEc2.Reservation) ([]*awsEc2.Instance, error) {
	if len(reservations) == 0 {
		return nil, errors.New("no instances provided")
	}

	var instances []*awsEc2.Instance
	for _, reservation := range reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
}

func (inst *EC2) GetRunningWorkerInstanceId() (string, error) {
	instances, err := inst.GetRunningInstances()
	if err != nil {
		return "", err
	}

	var instance *awsEc2.Instance
	for _, currentInstance := range instances {
		for _, securityGroup := range currentInstance.SecurityGroups {
			matched, err := filepath.Match("worker_group*", *securityGroup.GroupName)
			if err != nil {
				return "", err
			}
			if matched {
				instance = currentInstance
			}
		}
	}

	if instance == nil {
		return "", errors.New("unable to find a running EC2 worker instance")
	}

	return *instance.InstanceId, nil
}
