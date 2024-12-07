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

func NewDynamoDB(cfg *config.DynamoDB) (*DynamoDB, error) {
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

func dynamoDBClient(cfg *config.DynamoDB) (*dynamodb.Client, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfg.AwsRegion),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.AwsDynamoDBAccessKey,
				SecretAccessKey: cfg.AwsDynamoDBSecretey,
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}
