package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	dynamoClient *dynamodb.DynamoDB
	snsClient    *sns.SNS
	tableName    = "PlantMoistureData"
	topicArn     = os.Getenv("SNS_TOPIC_ARN")
)

type MoistureData struct {
	PlantID      string  `json:"plantId"`
	Timestamp    int64   `json:"timestamp"`
	MoistureLevel float64 `json:"moistureLevel"`
	Status       string  `json:"status"`
}

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
	snsClient = sns.New(sess)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var data MoistureData
	err := json.Unmarshal([]byte(event.Body), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	err = saveToDynamoDB(data)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Check moisture level and send notification if too low
	if data.MoistureLevel < 30 {
		err = sendNotification(data)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Data processed"}, nil
}

func saveToDynamoDB(data MoistureData) error {
	item := map[string]*dynamodb.AttributeValue{
		"PlantID": {
			S: aws.String(data.PlantID),
		},
		"Timestamp": {
			N: aws.String(fmt.Sprintf("%d", data.Timestamp)),
		},
		"MoistureLevel": {
			N: aws.String(fmt.Sprintf("%f", data.MoistureLevel)),
		},
		"Status": {
			S: aws.String(data.Status),
		},
	}

	_, err := dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func sendNotification(data MoistureData) error {
	message := fmt.Sprintf("Plant %s needs water! Moisture level: %.2f", data.PlantID, data.MoistureLevel)
	_, err := snsClient.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
	})
	return err
}

func main() {
	lambda.Start(handler)
}