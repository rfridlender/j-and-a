package models

type PersonMetadataPayload struct {
	GivenName  string
	FamilyName string
}

func (p *PersonMetadataPayload) Item(modelIdentifiers *ModelIdentifiers, version int, latestVersion int, createdAt string, createdBy string) ModelItem {
	return &PersonMetadataItem{
		GivenName:     p.GivenName,
		FamilyName:    p.FamilyName,
		PK:            EncodePartitionKey(ModelTypePerson, modelIdentifiers.PartitionId),
		SK:            EncodeSortKey(version, ModelTypePersonMetadata, modelIdentifiers.SortId),
		EntityType:    ModelTypePersonMetadata,
		LatestVersion: latestVersion,
		CreatedAt:     createdAt,
		CreatedBy:     createdBy,
		DeletedAt:     "",
		DeletedBy:     "",
	}
}

type PersonMetadataItem struct {
	GivenName     string
	FamilyName    string
	PK            string
	SK            string
	EntityType    string
	LatestVersion int `dynamodbav:",omitempty"`
	CreatedAt     string
	CreatedBy     string
	DeletedAt     string `dynamodbav:",omitempty"`
	DeletedBy     string `dynamodbav:",omitempty"`
}

func (i *PersonMetadataItem) Data() (ModelData, error) {
	_, personId, err := DecodePartitionKey(i.PK)
	if err != nil {
		return nil, err
	}

	return &PersonMetadataData{
		GivenName:  i.GivenName,
		FamilyName: i.FamilyName,
		PersonId:   personId,
		CreatedAt:  i.CreatedAt,
		CreatedBy:  i.CreatedBy,
		DeletedAt:  i.DeletedAt,
		DeletedBy:  i.DeletedBy,
	}, nil
}

type PersonMetadataData struct {
	GivenName  string
	FamilyName string
	PersonId   string
	CreatedAt  string
	CreatedBy  string
	DeletedAt  string
	DeletedBy  string
}
