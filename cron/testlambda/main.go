package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/silinternational/serverless-google-groups-sync/lib/syncgroups"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/silinternational/serverless-google-groups-sync"
)

const ApplicationConfigurationFile = "/tmp/config.json"

type RuntimeConfig struct {
	S3Bucket       string
	ConfigFilename string
	ConfigPath     string
	AppConfig      domain.AppConfig
}

func (t *RuntimeConfig) setRequired() error {
	errMsg := "Error: required value missing for environment variable %s"

	envKey := "S3_BUCKET"
	value := domain.GetEnv(envKey, "")
	if value == "" {
		return fmt.Errorf(errMsg, envKey)
	}
	t.S3Bucket = value

	return nil
}

func (t *RuntimeConfig) setDefaults() {
	if t.ConfigFilename == "" {
		t.ConfigFilename = "config.json"
	}
	if t.ConfigPath == "" {
		t.ConfigPath = "/tmp"
	}
}

func downloadConfigFromS3(config RuntimeConfig) error {

	bucketItem := s3.GetObjectInput{
		Bucket: aws.String(config.S3Bucket),
		Key:    aws.String(config.ConfigFilename),
	}

	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	s3Object, err := svc.GetObject(&bucketItem)
	if err != nil {
		log.Println("unable to get config file from S3 ... ")
		return err
	}

	body := s3Object.Body
	bodyBuf, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println("unable to read config file from S3:", err.Error())
		return err
	}

	err = ioutil.WriteFile(ApplicationConfigurationFile, bodyBuf, 0644)
	if err != nil {
		log.Println("unable to write config to disk:", err.Error())
		return err
	}

	return nil
}

func handler(runtimeConfig RuntimeConfig) error {
	log.Println("Starting TestLambda.")

	err := runtimeConfig.setRequired()
	if err != nil {
		log.Println("unable to set required runtimeConfig parameters, error:", err.Error())
		return err
	}

	runtimeConfig.setDefaults()

	err = downloadConfigFromS3(runtimeConfig)
	if err != nil {
		log.Println("unable to download runtimeConfig file from s3, error:", err.Error())
		return err
	}

	appConfig, err := domain.LoadAppConfig(ApplicationConfigurationFile)
	if err != nil {
		log.Println("unable to load app runtimeConfig, error:", err.Error())
		return err
	}

	err = syncgroups.SyncGroups(appConfig)
	if err != nil {
		log.Println("error running sync groups, error:", err.Error())
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
