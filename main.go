package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

var (
	svc        *dynamodb.DynamoDB
	err        error
	itemResult *dynamodb.GetItemOutput
)

type RequestStruct struct {
	Reg string `json:"reg"`
}

type ResponseStruct struct {
	Make  string `json:"make"`
	Model string `json:"model"`
}

func HandleRequest(ctx context.Context, request RequestStruct) (ResponseStruct, error) {
	log.Printf("Setting up request with reg: %s\n", request.Reg)
	itemInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"reg": {
				S: aws.String(request.Reg),
			},
		},
		TableName: aws.String("MOT"),
	}

	log.Println("Calling DDB getitem")
	itemResult, err = svc.GetItem(itemInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Panicln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Errorln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Errorln(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Errorln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Errorln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Errorln(err.Error())
		}
		return ResponseStruct{}, err
	}

	log.Println(itemResult)

	return ResponseStruct{
		Make: request.Reg,
	}, nil
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Println("Running init session")

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc = dynamodb.New(sess)
}

func main() {
	lambda.Start(HandleRequest)
}
