package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"

	"j-and-a/internal/models"
)

type Repository struct {
	Client    *dynamodb.Client
	TableName string
	IndexName string
}

func (r *Repository) DeleteByPartitionIdAndSortId(ctx context.Context, modelIdentifiers *models.ModelIdentifiers) error {
	getItemOutput, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.EncodePartitionKey(modelIdentifiers.PartitionType, modelIdentifiers.PartitionId)},
			"SK": &types.AttributeValueMemberS{Value: models.EncodeSortKey(0, modelIdentifiers.SortType, modelIdentifiers.SortId)},
		},
		ProjectionExpression: aws.String("LatestVersion"),
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

	unixMilli, ok := ctx.Value("requestedAt").(int64)
	if !ok {
		return errors.New("failed to parse requested at within context")
	}
	deletedAt := time.UnixMilli(unixMilli).Format(time.RFC3339)
	deletedBy, ok := ctx.Value("requestedBy").(string)
	if !ok {
		return errors.New("missing requested by within context")
	}

	_, err = r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Update: &types.Update{
				TableName: &r.TableName,
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: models.EncodePartitionKey(modelIdentifiers.PartitionType, modelIdentifiers.PartitionId)},
					"SK": &types.AttributeValueMemberS{Value: models.EncodeSortKey(0, modelIdentifiers.SortType, modelIdentifiers.SortId)},
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
					"PK": &types.AttributeValueMemberS{Value: models.EncodePartitionKey(modelIdentifiers.PartitionType, modelIdentifiers.PartitionId)},
					"SK": &types.AttributeValueMemberS{Value: models.EncodeSortKey(latestVersion, modelIdentifiers.SortType, modelIdentifiers.SortId)},
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

func (r *Repository) GetByPartitionId(ctx context.Context, modelIdentifiers *models.ModelIdentifiers, modelItem models.ModelItem) ([]models.ModelData, error) {
	queryOutput, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK, :SK)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: models.EncodePartitionKey(modelIdentifiers.PartitionType, modelIdentifiers.PartitionId)},
			":SK": &types.AttributeValueMemberS{Value: models.EncodeAnonymousSortKey(0, modelIdentifiers.SortType)},
		},
	})
	if err != nil {
		return nil, err
	}

	datas := make([]models.ModelData, queryOutput.Count)
	for idx, queryOutputItem := range queryOutput.Items {
		modelItem = modelItem.New()
		err = attributevalue.UnmarshalMap(queryOutputItem, modelItem)
		if err != nil {
			return nil, err
		}

		data, err := modelItem.Data()
		if err != nil {
			return nil, err
		}

		datas[idx] = data
	}

	return datas, nil
}

func (r *Repository) GetByPartitionIdAndSortId(ctx context.Context, modelIdentifiers *models.ModelIdentifiers, modelItem models.ModelItem) (models.ModelData, error) {
	getItemOutput, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.EncodePartitionKey(modelIdentifiers.PartitionType, modelIdentifiers.PartitionId)},
			"SK": &types.AttributeValueMemberS{Value: models.EncodeSortKey(0, modelIdentifiers.SortType, modelIdentifiers.SortId)},
		},
	})
	if err != nil {
		return nil, err
	}

	if getItemOutput.Item == nil {
		return nil, errors.New("item not found")
	}

	err = attributevalue.UnmarshalMap(getItemOutput.Item, modelItem)
	if err != nil {
		return nil, err
	}

	data, err := modelItem.Data()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) GetBySortType(ctx context.Context, modelIdentifiers *models.ModelIdentifiers, modelItem models.ModelItem) ([]models.ModelData, error) {
	queryOutput, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		IndexName:              aws.String(r.IndexName),
		KeyConditionExpression: aws.String("ModelType = :ModelType AND begins_with(SK, :SK)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ModelType": &types.AttributeValueMemberS{Value: string(modelIdentifiers.SortType)},
			":SK":        &types.AttributeValueMemberS{Value: models.EncodeAnonymousSortKey(0, modelIdentifiers.SortType)},
		},
	})
	if err != nil {
		return nil, err
	}

	datas := make([]models.ModelData, queryOutput.Count)
	for idx, queryOutputItem := range queryOutput.Items {
		modelItem = modelItem.New()
		err = attributevalue.UnmarshalMap(queryOutputItem, modelItem)
		if err != nil {
			return nil, err
		}

		data, err := modelItem.Data()
		if err != nil {
			return nil, err
		}

		datas[idx] = data
	}

	return datas, nil
}

func (r *Repository) PutByPartitionIdAndSortId(ctx context.Context, modelIdentifiers *models.ModelIdentifiers, modelPayload models.ModelPayload) error {
	getItemOutput, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.EncodePartitionKey(modelIdentifiers.PartitionType, modelIdentifiers.PartitionId)},
			"SK": &types.AttributeValueMemberS{Value: models.EncodeSortKey(0, modelIdentifiers.SortType, modelIdentifiers.SortId)},
		},
		ProjectionExpression: aws.String("LatestVersion"),
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

	unixMilli, ok := ctx.Value("requestedAt").(int64)
	if !ok {
		return errors.New("failed to parse requested at within context")
	}
	createdAt := time.UnixMilli(unixMilli).Format(time.RFC3339)
	createdBy, ok := ctx.Value("requestedBy").(string)
	if !ok {
		return errors.New("failed to parse requested by within context")
	}

	rootItem, err := attributevalue.MarshalMap(modelPayload.Item(modelIdentifiers, 0, latestVersion+1, createdAt, createdBy))
	if err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(modelPayload.Item(modelIdentifiers, latestVersion+1, 0, createdAt, createdBy))
	if err != nil {
		return err
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
