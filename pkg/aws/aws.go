package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsEc2 "github.com/aws/aws-sdk-go/service/ec2"
	awsElasticSearch "github.com/aws/aws-sdk-go/service/elasticsearchservice"
	awsRds "github.com/aws/aws-sdk-go/service/rds"
	"github.com/ryannel/hippo/pkg/aws/ec2"
	"github.com/ryannel/hippo/pkg/aws/elasticSearch"
	"github.com/ryannel/hippo/pkg/aws/rds"
	"github.com/ryannel/hippo/pkg/logger"
	"github.com/ryannel/hippo/pkg/util"
)

type Aws struct {
	EC2 ec2.EC2
	ElasticSearch elasticSearch.ElasticSearch
	RDS rds.RDS
}

func New(region string) (Aws, error){
	sess, _ := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:   aws.String(region),
		},
	})

	awsInstance := Aws {
		EC2: ec2.EC2{Sdk :awsEc2.New(sess)},
		ElasticSearch: elasticSearch.ElasticSearch{Sdk: awsElasticSearch.New(sess)},
		RDS: rds.RDS{Sdk: awsRds.New(sess)},
	}

	return awsInstance, nil
}

func (aws *Aws) Login(profile string) (string, error) {
	profileExtension := ""
	if profile != "" {
		profileExtension = " --profile " + profile
	}
	command := "aws-azure-login.cmd --no-prompt" + profileExtension
	logger.Command("Assuming AWS role: " + command)
	return util.ExecStringCommand(command)
}