package aws

import (
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
	awsCfg, err := awsClient(cfgAWS)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}
