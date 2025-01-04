package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/idoyudha/eshop-product/config"
)

type S3Service struct {
	Client        *s3.Client
	ProductBucket string
	CDNDomain     string
}

func NewS3(cfg *config.AWS) (*S3Service, error) {
	s3Service := &S3Service{
		ProductBucket: cfg.ProductBucket,
		CDNDomain:     cfg.CdnDomain,
	}

	client, err := s3Client(cfg)
	if err != nil {
		return nil, err
	}

	s3Service.Client = client
	return s3Service, nil
}

func s3Client(cfg *config.AWS) (*s3.Client, error) {
	awsCfg, err := awsClient(cfg)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg), nil
}
