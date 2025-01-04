package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/idoyudha/eshop-product/config"
)

func awsClient(cfgAWS *config.AWS) (aws.Config, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfgAWS.AwsRegion),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfgAWS.AwsAccessKey,
				SecretAccessKey: cfgAWS.AwsSecretey,
			},
		}),
	)
	if err != nil {
		return aws.Config{}, err
	}
	return awsCfg, nil
}
