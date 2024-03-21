package workflow

import (
	"context"
	"example/src/config"
	"fmt"
	"log"
	"os"

	"github.com/akyoto/uuid"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/canhlinh/hlsdl"
	"go.temporal.io/sdk/activity"
)

/**
 * Sample activities used by file processing sample workflow.
 */

type Activities struct {
	S3 *s3.Client
}

func (a *Activities) DownloadFileActivity(ctx context.Context, fileName string, linkMediaFile string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Downloading file...", "FileID", fileName)
	hlsDL := hlsdl.New(linkMediaFile, nil, "download", uuid.New().String()+".mp4", 64, true)

	filepath, err := hlsDL.Download()
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func (a *Activities) UploadActivity(ctx context.Context, fileNameLocal string, key string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("upload to s3", "FileName", fileNameLocal)
	file, err := os.Open(fileNameLocal)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}

	defer file.Close()
	endpoint := fmt.Sprintf("http://%s:%+v", config.GetConfiguration().S3.Host, config.GetConfiguration().S3.Port)

	accessKey := config.GetConfiguration().S3.AccessKey
	secretKey := config.GetConfiguration().S3.SecretKey
	region := config.GetConfiguration().S3.Region
	log.Printf("%s %s", accessKey, secretKey)
	//disableSSL := true
	//staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
	//	return aws.Endpoint{
	//		PartitionID:       "aws",
	//		URL:               endpoint, // or where ever you ran minio
	//		SigningRegion:     region,
	//		HostnameImmutable: true,
	//	}, nil
	//})

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			HostnameImmutable: true,
			PartitionID:       "aws",
			SigningName:       "",
			SigningRegion:     region,
			SigningMethod:     "",
			Source:            0,
		}, nil
	})
	cfg := aws.Config{
		Region:                      region,
		Credentials:                 credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		BearerAuthTokenProvider:     nil,
		HTTPClient:                  nil,
		EndpointResolver:            nil,
		EndpointResolverWithOptions: resolver,
		RetryMaxAttempts:            0,
		RetryMode:                   "",
		Retryer:                     nil,
		ConfigSources:               nil,
		APIOptions:                  nil,
		Logger:                      nil,
		ClientLogMode:               0,
		DefaultsMode:                "",
		RuntimeEnvironment:          aws.RuntimeEnvironment{},
	}
	s3Connect := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	uploader := manager.NewUploader(s3Connect)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.GetConfiguration().S3.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		fmt.Println("S3:", err)
		return "", err
	}
	defer func() { _ = os.Remove(fileNameLocal) }() // cleanup temp file

	// read downloaded file

	return "", nil
}

func NewActivitie() *Activities {

	// log.Printf("init %+v", s3)
	return &Activities{}

}
