package uploadService

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/repository"
)

type UploadService interface {
	HandlerUploadFile(c *fiber.Ctx) error
}

type uploadService struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository repository.EncryptorRepository
	s3Client            *s3.Client
}

func ProvideUploadService(db *gorm.DB, config config.AppConfig) UploadService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)

	awsCfg := aws.Config{
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(config.S3Config.S3AccessKey, config.S3Config.S3SecretKey, ""),
		),
		Region: "auto",
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:           config.S3Config.S3Endpoint,
					SigningRegion: "auto",
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		}),
	}

	s3Client := s3.NewFromConfig(awsCfg)
	return &uploadService{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
		s3Client:            s3Client,
	}
}
