package models

import (
	utils "j-and-a/pkg"
)

const LOG_PARTITION_TYPE = "JOB"
const LOG_SORT_TYPE = "LOG"

type LogRequest struct {
	PersonId string
	Hours    int
	JobId    string
	LogId    string
}

func (r *LogRequest) ToItem(version int, latestVersion int, createdAt string, createdBy string) *LogItem {
	return &LogItem{
		PersonId:      r.PersonId,
		Hours:         r.Hours,
		PK:            utils.EncodePartitionKey(LOG_PARTITION_TYPE, r.JobId),
		SK:            utils.EncodeSortKey(version, LOG_SORT_TYPE, r.LogId),
		EntityType:    LOG_SORT_TYPE,
		LatestVersion: latestVersion,
		CreatedAt:     createdAt,
		CreatedBy:     createdBy,
		DeletedAt:     "",
		DeletedBy:     "",
	}
}

type LogItem struct {
	PersonId      string
	Hours         int
	PK            string
	SK            string
	EntityType    string
	LatestVersion int `dynamodbav:",omitempty"`
	CreatedAt     string
	CreatedBy     string
	DeletedAt     string `dynamodbav:",omitempty"`
	DeletedBy     string `dynamodbav:",omitempty"`
}

func (i *LogItem) ToData() (*LogData, error) {
	_, jobId, err := utils.DecodePartitionKey(i.PK)
	if err != nil {
		return nil, err
	}

	_, _, logId, err := utils.DecodeSortKey(i.SK)
	if err != nil {
		return nil, err
	}
	return &LogData{
		PersonId:  i.PersonId,
		Hours:     i.Hours,
		JobId:     jobId,
		LogId:     logId,
		CreatedAt: i.CreatedAt,
		CreatedBy: i.CreatedBy,
		DeletedAt: i.DeletedAt,
		DeletedBy: i.DeletedBy,
	}, nil
}

type LogData struct {
	PersonId  string
	Hours     int
	JobId     string
	LogId     string
	CreatedAt string
	CreatedBy string
	DeletedAt string
	DeletedBy string
}
