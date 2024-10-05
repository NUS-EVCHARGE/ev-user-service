package thirdparty
// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"log"
)

type SecretManagerController interface {
	GetSecretFromManager(secretName string) *string
	GetPrivateSignKey() *string
	SetPrivateSignKey(privateKey string)
}

type SecretManagerControllerImpl struct {
}

var (
	SecretManagerControllerObj SecretManagerController
	PrivateKeySigning          *string
)

func NewSecretManagerController() {
	SecretManagerControllerObj = &SecretManagerControllerImpl{}
}

func (s SecretManagerControllerImpl) GetSecretFromManager(secretName string) *string {
	region := "ap-southeast-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString = *result.SecretString

	return &secretString
}

func (s SecretManagerControllerImpl) GetPrivateSignKey() *string {
	return PrivateKeySigning
}

func (s SecretManagerControllerImpl) SetPrivateSignKey(privateKey string) {
	PrivateKeySigning = &privateKey
}