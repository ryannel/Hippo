package aws

import "testing"

func Test_Live(t *testing.T) {
	aws, _ := Connect("", "eu-west-1")
	rds := aws.RDS
	instances, _ := rds.GetInstances()

	endpoint, _ := rds.GetEndpoint(instances[0])

	println(endpoint.Address)
}