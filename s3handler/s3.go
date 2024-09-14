package s3handler

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/config"
)

var (
	AwsSession *session.Session
	S3Client   *s3.S3
)

// InitializeS3Session initializes an AWS session
func InitializeS3Session(awsConfig config.AWSConfig) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsConfig.Region),
		Credentials: credentials.NewStaticCredentials(awsConfig.AccessKey, awsConfig.SecretKey, ""),
	})
	if err != nil {
		log.Fatalf("failed to initialize AWS session: %v", err)
	}
	AwsSession = sess
	S3Client = s3.New(sess)
	log.Println("AWS session and S3 client initialized successfully")
}
