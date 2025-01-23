package models

type LogPayload struct {
	PersonId string  `json:"personId"`
	Hours    float64 `json:"hours"`
}

func (p *LogPayload) Item(modelIdentifiers *ModelIdentifiers, version int, latestVersion int, createdAt string, createdBy string) ModelItem {
	return &LogItem{
		PersonId:      p.PersonId,
		Hours:         p.Hours,
		PK:            EncodePartitionKey(ModelTypeJob, modelIdentifiers.PartitionId),
		SK:            EncodeSortKey(version, ModelTypeLog, modelIdentifiers.SortId),
		ModelType:     ModelTypeLog,
		LatestVersion: latestVersion,
		CreatedAt:     createdAt,
		CreatedBy:     createdBy,
		DeletedAt:     "",
		DeletedBy:     "",
	}
}

type LogItem struct {
	PersonId      string
	Hours         float64
	PK            string
	SK            string
	ModelType     string
	LatestVersion int `dynamodbav:",omitempty"`
	CreatedAt     string
	CreatedBy     string
	DeletedAt     string `dynamodbav:",omitempty"`
	DeletedBy     string `dynamodbav:",omitempty"`
}

func (i *LogItem) Data() (ModelData, error) {
	_, partitionId, err := DecodePartitionKey(i.PK)
	if err != nil {
		return nil, err
	}

	_, _, sortId, err := DecodeSortKey(i.SK)
	if err != nil {
		return nil, err
	}

	return &LogData{
		PersonId:  i.PersonId,
		Hours:     i.Hours,
		JobId:     partitionId,
		LogId:     sortId,
		CreatedAt: i.CreatedAt,
		CreatedBy: i.CreatedBy,
		DeletedAt: i.DeletedAt,
		DeletedBy: i.DeletedBy,
	}, nil
}

type LogData struct {
	PersonId  string  `json:"personId"`
	Hours     float64 `json:"hours"`
	JobId     string  `json:"jobId"`
	LogId     string  `json:"logId"`
	CreatedAt string  `json:"createdAt"`
	CreatedBy string  `json:"createdBy"`
	DeletedAt string  `json:"deletedAt"`
	DeletedBy string  `json:"deletedBy"`
}
