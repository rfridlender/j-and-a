package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const NUMBER_OF_PARTITION_KEY_PARTS = 2
const NUMBER_OF_SORT_KEY_PARTS = 3

const SORT_KEY_VERSION_PREFIX = "V"

type ModelType string

const (
	ModelTypePerson         = "Person"
	ModelTypePersonMetadata = "PersonMetadata"
)

type ModelIdentifiers struct {
	PartitionType ModelType
	PartitionId   string
	SortType      ModelType
	SortId        string
}

type ModelPayload interface {
	Item(modelIdentifiers *ModelIdentifiers, version int, latestVersion int, createdAt string, createdBy string) ModelItem
}

type ModelItem interface {
	Data() (ModelData, error)
}

type ModelData interface{}

func EncodePartitionKey(partitionType ModelType, partitionId string) string {
	return fmt.Sprintf("%s#%s", partitionType, partitionId)
}

func EncodeSortKey(version int, sortType ModelType, sortId string) string {
	return fmt.Sprintf("%s%d#%s#%s", SORT_KEY_VERSION_PREFIX, version, sortType, sortId)
}

func EncodeAnonymousSortKey(version int, sortType ModelType) string {
	return fmt.Sprintf("%s%d#%s#", SORT_KEY_VERSION_PREFIX, version, sortType)
}

func DecodePartitionKey(partitionKey string) (ModelType, string, error) {
	parts := strings.Split(partitionKey, "#")
	if len(parts) != NUMBER_OF_PARTITION_KEY_PARTS {
		return "", "", errors.New("invalid partition key")
	}
	return ModelType(parts[0]), parts[1], nil
}

func DecodeSortKey(sortKey string) (int, ModelType, string, error) {
	parts := strings.Split(sortKey, "#")
	if len(parts) != NUMBER_OF_SORT_KEY_PARTS {
		return 0, "", "", errors.New("invalid sort key")
	}
	versionString := strings.TrimPrefix(parts[0], SORT_KEY_VERSION_PREFIX)
	if versionString == parts[0] {
		return 0, "", "", errors.New("invalid version")
	}
	version, err := strconv.Atoi(versionString)
	if err != nil {
		return 0, "", "", err
	}
	return version, ModelType(parts[1]), parts[2], nil
}
