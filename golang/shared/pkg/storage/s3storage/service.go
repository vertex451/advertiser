package s3storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"io"
	"log"
)

const bucketName = "stage-tg-bot-public-bucket"

type Service struct {
	cfg        aws.Config
	bucketName string
	s3Client   *s3.Client
}

func New() *Service {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		cfg:        cfg,
		bucketName: bucketName,
		s3Client:   s3.NewFromConfig(cfg),
	}
}

func (s *Service) Store(name string, body io.Reader) (string, error) {
	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    aws.String(name),
		Body:   body,
	})
	if err != nil {
		zap.L().Error("failed to put object to bucket", zap.Error(err))
		return "", err
	}

	return fmt.Sprintf("s3://%s.s3.%s.amazonaws.com/%s", bucketName, s.cfg.Region, name), nil
}

func (s *Service) Load() {

}
