package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
	"j-and-a/internal/services"
)

var (
	client    *dynamodb.Client
	tableName string
	indexName string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	client = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("DYNAMO_DB_TABLE_NAME")
	indexName = os.Getenv("DYNAMO_DB_INDEX_NAME")
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
    if request.IsBase64Encoded {
        decodedRequestBody, err := base64.StdEncoding.DecodeString(request.Body)
        if err != nil {
            log.Fatal(err)
        }
        request.Body = string(decodedRequestBody)
    }

    jsonRequest, err := json.Marshal(request)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("request %s", string(jsonRequest))

	ctx = context.WithValue(ctx, "requestedBy", request.RequestContext.Authorizer.JWT.Claims["sub"])

	repository := &repositories.LogRepository{
		Client:    client,
		TableName: tableName,
		IndexName: indexName,
	}

	service := &services.LogService{Repository: repository}

	var data interface{}
	switch request.RouteKey {
	case "GET /jobs/{jobId}/logs":
		data, err = service.GetAllLogs(ctx)
	case "GET /jobs/{jobId}/logs/{logId}":
		jobId := request.PathParameters["jobId"]
		logId := request.PathParameters["logId"]
		data, err = service.GetLog(ctx, jobId, logId)
	case "PUT /jobs/{jobId}/logs/{logId}":
		newRequest := new(models.LogRequest)
		err = json.Unmarshal([]byte(request.Body), newRequest)
		if err != nil {
			break
		}
        newRequest.JobId = request.PathParameters["jobId"]
        newRequest.LogId = request.PathParameters["logId"]
		err = service.PutLog(ctx, newRequest)
	case "DELETE /jobs/{jobId}/logs/{logId}":
		jobId := request.PathParameters["jobId"]
		logId := request.PathParameters["logId"]
		err = service.DeleteLog(ctx, jobId, logId)
	default:
		err = errors.New("unsupported request.RouteKey")
	}
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonData),
	}, nil
}

func main() {
	lambda.Start(handler)
}
