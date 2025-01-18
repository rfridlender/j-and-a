package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"

	"j-and-a/internal/models"
	utils "j-and-a/pkg"
)

type LogRepository struct {
	Client    *dynamodb.Client
	TableName string
	IndexName string
}

func (r *LogRepository) Delete(ctx context.Context, partitionId string, sortId string) error {
	getItemOutput, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: utils.EncodePartitionKey(models.LOG_PARTITION_TYPE, partitionId)},
			"SK": &types.AttributeValueMemberS{Value: utils.EncodeSortKey(0, models.LOG_SORT_TYPE, sortId)},
		},
	})
	if err != nil {
		return err
	}

	latestVersion := 0
	if lastedVersionAttributeValue, ok := getItemOutput.Item["LatestVersion"]; ok {
		err = attributevalue.Unmarshal(lastedVersionAttributeValue, &latestVersion)
		if err != nil {
			return err
		}
	}

	deletedAt := time.Now().Format(time.RFC3339)
	deletedBy, ok := ctx.Value("requestedBy").(string)
	if !ok {
		return errors.New("missing requested by within context")
	}

	_, err = r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Update: &types.Update{
				TableName: &r.TableName,
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: utils.EncodePartitionKey(models.LOG_PARTITION_TYPE, partitionId)},
					"SK": &types.AttributeValueMemberS{Value: utils.EncodeSortKey(0, models.LOG_SORT_TYPE, sortId)},
				},
				UpdateExpression: aws.String("SET DeletedAt = :DeletedAt, DeletedBy = :DeletedBy"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":DeletedAt": &types.AttributeValueMemberS{Value: deletedAt},
					":DeletedBy": &types.AttributeValueMemberS{Value: deletedBy},
				},
				ConditionExpression: aws.String("attribute_not_exists(deletedAt)"),
			}},
			{Update: &types.Update{
				TableName: &r.TableName,
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: utils.EncodePartitionKey(models.LOG_PARTITION_TYPE, partitionId)},
					"SK": &types.AttributeValueMemberS{Value: utils.EncodeSortKey(latestVersion, models.LOG_SORT_TYPE, sortId)},
				},
				UpdateExpression: aws.String("SET DeletedAt = :DeletedAt, DeletedBy = :DeletedBy"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":DeletedAt": &types.AttributeValueMemberS{Value: deletedAt},
					":DeletedBy": &types.AttributeValueMemberS{Value: deletedBy},
				},
				ConditionExpression: aws.String("attribute_not_exists(deletedAt)"),
			}},
		},
	})

	return err
}

func (r *LogRepository) GetAll(ctx context.Context) ([]models.LogData, error) {
	queryOutput, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		IndexName:              aws.String(r.IndexName),
		KeyConditionExpression: aws.String("EntityType = :EntityType AND begins_with(SK, :SK)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":EntityType": &types.AttributeValueMemberS{Value: models.LOG_SORT_TYPE},
			":SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d#%s#", utils.SORT_KEY_VERSION_PREFIX, 0, models.LOG_SORT_TYPE)},
		},
	})
	if err != nil {
		return nil, err
	}

	datas := make([]models.LogData, queryOutput.Count)
	for idx, queryOutputItem := range queryOutput.Items {
		item := new(models.LogItem)
		err = attributevalue.UnmarshalMap(queryOutputItem, item)
		if err != nil {
			return nil, err
		}

		data, err := item.ToData()
		if err != nil {
			return nil, err
		}

		datas[idx] = *data
	}

	return datas, nil
}

func (r *LogRepository) Get(ctx context.Context, partitionId string, sortId string) (*models.LogData, error) {
	getItemOutput, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: utils.EncodePartitionKey(models.LOG_PARTITION_TYPE, partitionId)},
			"SK": &types.AttributeValueMemberS{Value: utils.EncodeSortKey(0, models.LOG_SORT_TYPE, sortId)},
		},
	})
	if err != nil {
		return nil, err
	}

	if getItemOutput.Item == nil {
		return nil, errors.New("item not found")
	}

	item := new(models.LogItem)
	err = attributevalue.UnmarshalMap(getItemOutput.Item, item)
	if err != nil {
		return nil, err
	}

	data, err := item.ToData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *LogRepository) Put(ctx context.Context, request *models.LogRequest) error {
	getItemOutput, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: utils.EncodePartitionKey(models.LOG_PARTITION_TYPE, request.JobId)},
			"SK": &types.AttributeValueMemberS{Value: utils.EncodeSortKey(0, models.LOG_SORT_TYPE, request.LogId)},
		},
	})
	if err != nil {
		return err
	}

	latestVersion := 0
	if lastedVersionAttributeValue, ok := getItemOutput.Item["LatestVersion"]; ok {
		err = attributevalue.Unmarshal(lastedVersionAttributeValue, &latestVersion)
		if err != nil {
			return err
		}
	}

	createdAt := time.Now().Format(time.RFC3339)
	createdBy, ok := ctx.Value("requestedBy").(string)
	if !ok {
		return errors.New("missing requested by within context")
	}

	rootItem, err := attributevalue.MarshalMap(request.ToItem(0, latestVersion+1, createdAt, createdBy))
	if err != nil {
		log.Fatal(err)
	}

	item, err := attributevalue.MarshalMap(request.ToItem(latestVersion+1, 0, createdAt, createdBy))
	if err != nil {
		log.Fatal(err)
	}

	_, err = r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Put: &types.Put{
				TableName: &r.TableName,
				Item:      rootItem,
			}},
			{Put: &types.Put{
				TableName: &r.TableName,
				Item:      item,
			}},
		},
	})

	return err
}
