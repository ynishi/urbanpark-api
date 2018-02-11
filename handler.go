package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/ynishi/urbanparkjp"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	regionName = "ap-northeast-1"
	tableName  = "parks"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {

	id, ok := request.QueryStringParameters["id"]
	if !ok {
		response = events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
		return response, err
	}
	log.Printf("id: %s\n", id)

	if len(os.Getenv("DYNAMO_REGION")) > 0 {
		regionName = os.Getenv("DYNAMO_REGION")
	}
	log.Printf("dynamo region: %s\n", regionName)
	if len(os.Getenv("DYNAMO_TABLE")) > 0 {
		tableName = os.Getenv("DYNAMO_TABLE")
	}
	log.Printf("dynamo table: %s\n", tableName)

	sess, err := session.NewSession(&aws.Config{Region: aws.String(regionName)})
	if err != nil {
		response = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		return response, err
	}
	db := dynamo.New(sess)
	table := db.Table(tableName)

	var result urbanparkjp.Park

	err = table.Get("id", id).One(&result)
	if err == dynamo.ErrNotFound {
		response = events.APIGatewayProxyResponse{
			Body:            `{"message": "not found"}"`,
			StatusCode:      http.StatusOK,
			IsBase64Encoded: true,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return response, nil
	}
	if err != nil {
		response = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		return response, err
	}
	log.Printf("park : %v\n", result)

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		response = events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	response = events.APIGatewayProxyResponse{
		Body:            string(jsonBytes),
		StatusCode:      http.StatusOK,
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
