package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/idoyudha/eshop-product/config"
)

const (
	_productTableName  = "eshop-products"
	_categoryTableName = "eshop-product-categories"
)

type DynamoDB struct {
	Client        *dynamodb.Client
	ProductTable  string
	CategoryTable string
}

func NewDynamoDB(cfg *config.AWS) (*DynamoDB, error) {
	dynamoDB := &DynamoDB{
		ProductTable:  _productTableName,
		CategoryTable: _categoryTableName,
	}

	client, err := dynamoDBClient(cfg)
	if err != nil {
		return nil, err
	}

	dynamoDB.Client = client
	return dynamoDB, nil
}

func dynamoDBClient(cfgAWS *config.AWS) (*dynamodb.Client, error) {
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
		return nil, err
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}
