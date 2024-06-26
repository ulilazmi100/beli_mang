package svc

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"beli_mang/responses"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// ImageSvc interface with the UploadImage method
type ImageSvc interface {
	UploadImage(*multipart.FileHeader) (string, error)
}

// imageSvc struct implementing the ImageSvc interface
type imageSvc struct {
	s3Client *s3.Client
	bucket   string
}

// NewImageSvc creates a new ImageSvc instance
func NewImageSvc(cfg aws.Config, bucket string) ImageSvc {
	s3Client := s3.NewFromConfig(cfg)
	return &imageSvc{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

// UploadImage handles the uploading of an image to S3
func (i *imageSvc) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	uuid := uuid.NewString()
	splitted := strings.Split(fileHeader.Filename, ".")
	fileExt := strings.ToLower(splitted[len(splitted)-1])
	fileName := fmt.Sprintf("%s.%s", uuid, fileExt)

	if fileExt != "jpg" && fileExt != "jpeg" {
		return "", responses.NewBadRequestError("file must be in JPEG format")
	}

	if fileHeader.Size < 10*1024 || fileHeader.Size > 2*1024*1024 {
		return "", responses.NewBadRequestError("file size must be between 10KB and 2MB")
	}

	if fileHeader.Header.Get("Content-Type") != "image/jpeg" {
		return "", responses.NewBadRequestError("file content type must be image/jpeg")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(fileName),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*2)
	defer cancel()

	_, err = i.s3Client.PutObject(ctx, input)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.ap-southeast-1.amazonaws.com/%s", i.bucket, fileName)

	return url, nil
}
