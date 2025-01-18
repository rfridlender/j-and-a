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

	repository := &repositories.PersonRepository{
		Client:    client,
		TableName: tableName,
		IndexName: indexName,
	}

	service := &services.PersonService{Repository: repository}

	var data any
	switch request.RouteKey {
	case "GET /persons":
		data, err = service.GetAllPersons(ctx)
	case "GET /persons/{personId}":
		personId := request.PathParameters["personId"]
		data, err = service.GetPerson(ctx, personId)
	case "PUT /persons/{personId}":
		newRequest := new(models.PersonRequest)
		err = json.Unmarshal([]byte(request.Body), newRequest)
		if err != nil {
			break
		}
		newRequest.PersonId = request.PathParameters["personId"]
		err = service.PutPerson(ctx, newRequest)
	case "DELETE /persons/{personId}":
		personId := request.PathParameters["personId"]
		err = service.DeletePerson(ctx, personId)
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
