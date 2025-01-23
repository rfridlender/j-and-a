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
			return nil, err
		}
		request.Body = string(decodedRequestBody)
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	log.Printf("request %s", string(jsonRequest))

	ctx = context.WithValue(ctx, "requestedBy", request.RequestContext.Authorizer.JWT.Claims["sub"])

	repository := &repositories.Repository{Client: client, TableName: tableName, IndexName: indexName}

	modelIdentifiers := &models.ModelIdentifiers{
		PartitionType: models.ModelType(request.PathParameters["PartitionType"]),
		PartitionId:   request.PathParameters["PartitionId"],
		SortType:      models.ModelType(request.PathParameters["SortType"]),
		SortId:        request.PathParameters["SortId"],
	}

	service, err := services.New(repository, modelIdentifiers)
	if err != nil {
		return nil, err
	}

	var data interface{}
	switch request.RouteKey {
	case "DELETE /{PartitionType}/{PartitionId}/{SortType}", "DELETE /{PartitionType}/{PartitionId}/{SortType}/{SortId}":
		err = service.DeleteByPartitionIdAndSortId(ctx)
	case "GET /{PartitionType}/{PartitionId}/{SortType}":
		data, err = service.GetByPartitionId(ctx)
	case "GET /{PartitionType}/{PartitionId}/{SortType}/{SortId}":
		data, err = service.GetByPartitionIdAndSortId(ctx)
	case "GET /{SortType}":
		data, err = service.GetBySortType(ctx)
	case "PUT /{PartitionType}/{PartitionId}/{SortType}", "PUT /{PartitionType}/{PartitionId}/{SortType}/{SortId}":
		err = service.PutByPartitionIdAndSortId(ctx, request.Body)
	default:
		err = errors.New("unsupported request.RouteKey")
	}

	if err != nil {
		log.Fatal(err)
	}

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       string(jsonData),
		}, nil
	} else {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
