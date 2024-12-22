package repo

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	awsService "github.com/idoyudha/eshop-product/pkg/aws"
)

type ProductS3Repo struct {
	*awsService.S3Service
}

func NewProductS3Repo(s *awsService.S3Service) *ProductS3Repo {
	return &ProductS3Repo{
		s,
	}
}

func (r *ProductS3Repo) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	filename := fmt.Sprintf("product/%s%s", uuid.New().String(), filepath.Ext(file.Filename))

	_, err = r.S3Service.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.ProductBucket),
		Key:         aws.String(filename),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	imageURL := fmt.Sprintf("https://%s/%s", r.CDNDomain, filename)

	return imageURL, nil
}
