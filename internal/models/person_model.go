package models

import (
	utils "j-and-a/pkg"
)

const PERSON_PARTITION_TYPE = "PERSON"
const PERSON_SORT_TYPE = "PERSON"

type PersonRequest struct {
	PersonId   string `json:"personId"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

func (r *PersonRequest) ToItem(version int, latestVersion int, createdAt string, createdBy string) *PersonItem {
	return &PersonItem{
		GivenName:     r.GivenName,
		FamilyName:    r.FamilyName,
		PK:            utils.EncodePartitionKey(PERSON_PARTITION_TYPE, r.PersonId),
		SK:            utils.EncodeSortKey(version, PERSON_SORT_TYPE, r.PersonId),
		EntityType:    PERSON_SORT_TYPE,
		LatestVersion: latestVersion,
		CreatedAt:     createdAt,
		CreatedBy:     createdBy,
		DeletedAt:     "",
		DeletedBy:     "",
	}
}

type PersonItem struct {
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

func (i *PersonItem) ToData() (*PersonData, error) {
	_, personId, err := utils.DecodePartitionKey(i.PK)
	if err != nil {
		return nil, err
	}

	return &PersonData{
		GivenName:  i.GivenName,
		FamilyName: i.FamilyName,
		PersonId:   personId,
		CreatedAt:  i.CreatedAt,
		CreatedBy:  i.CreatedBy,
		DeletedAt:  i.DeletedAt,
		DeletedBy:  i.DeletedBy,
	}, nil
}

type PersonData struct {
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	PersonId   string `json:"personId"`
	CreatedAt  string `json:"createdAt"`
	CreatedBy  string `json:"createdBy"`
	DeletedAt  string `json:"deletedAt"`
	DeletedBy  string `json:"deletedBy"`
}
