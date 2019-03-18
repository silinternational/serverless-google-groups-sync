package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)


type TestConfig struct {
	GroupsMapS3ARN string
	GroupsMapFileName string `json:"GroupsMapFileName"`
	AWSAccessKeyID string
	AWSSecretAccessKey string
}

func (t *TestConfig) setRequired() error {
	errMsg := "Error: required value missing for environment variable %s"


	envKey := "S3_BUCKET_FOR_GROUPS_MAP"
	value := os.Getenv(envKey)
	if value == "" {
		return fmt.Errorf(errMsg, envKey)
	}
	t.GroupsMapS3ARN = value

	return nil
}

func (t *TestConfig) setDefaults() {
	if t.GroupsMapFileName == "" {
		t.GroupsMapFileName = "groups_map.json"
	}
}


func handler(config TestConfig) error {
	log.Println("Starting TestLambda.")

	err := config.setRequired()

	if err != nil {
		return err
	}

	config.setDefaults()

	log.Println("Groups Map S3: ", config.GroupsMapS3ARN, " / ", config.GroupsMapFileName)

	bucketItem := s3.GetObjectInput{
		Bucket: aws.String(config.GroupsMapS3ARN),
		Key: aws.String(config.GroupsMapFileName),
	}

	sess := session.Must(session.NewSession())

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	svc := s3.New(sess)

	object, err := svc.GetObject(&bucketItem)

	if err != nil {
		log.Println("Unable to get Groups Map file from S3 ... ")
		return err
	}

	log.Println("Success: S3 file length: ", fmt.Sprintf("%v", object.ContentLength))


	return nil
}


func main() {
	lambda.Start(handler)
}
