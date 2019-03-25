package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/silinternational/serverless-google-groups-sync"
)

const GoogleCredsJsonFile = "/tmp/google-creds.json"

type TestConfig struct {
	GroupsMapS3ARN        string
	GroupsMapFileName     string `json:"GroupsMapFileName"`
	MemberSourceApiConfig domain.MemberSourceApiConfig
	// AWSAccessKeyID     string
	// AWSSecretAccessKey string
}

func (t *TestConfig) setRequired() error {
	errMsg := "Error: required value missing for environment variable %s"

	envKey := "S3_BUCKET_FOR_INPUT"
	value := domain.GetEnv(envKey, "")
	if value == "" {
		return fmt.Errorf(errMsg, envKey)
	}
	t.GroupsMapS3ARN = value

	return nil
}

func (t *TestConfig) setDefaults() {
	if t.GroupsMapFileName == "" {
		t.GroupsMapFileName = "groups-map.json"
	}
}

func saveGoogleCredsJsonFile(objectOutput *s3.GetObjectOutput) error {
	body := objectOutput.Body

	bodyBuf, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("unable to read Google Credentials file from S3: %s", err.Error())
	}

	err = ioutil.WriteFile(GoogleCredsJsonFile, bodyBuf, 0644)
	if err != nil {
		return fmt.Errorf("unable to write Google Credentials to disk: %s", err.Error())
	}

	// log.Println("Wrote following to disk: \n ", string(bodyBuf))
	return nil
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
		Key:    aws.String(config.GroupsMapFileName),
	}

	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	s3Object, err := svc.GetObject(&bucketItem)

	if err != nil {
		log.Println("Unable to get Groups Map file from S3 ... ")
		return err
	}

	// TODO Call a function that converts the Groups Map file from json into a golang map

	// TODO Call this with the actual google credentials json file (not the Groups Map file)
	err = saveGoogleCredsJsonFile(s3Object)

	if err != nil {
		return err
	}

	log.Println("Success: S3 file length: ", fmt.Sprintf("%d", &s3Object.ContentLength))

	return nil
}

func main() {
	lambda.Start(handler)
}
