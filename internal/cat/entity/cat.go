package cat

import (
	"time"

	"github.com/lib/pq"
)

type Cat struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Race        Race           `json:"race"`
	Sex         Sex            `json:"sex"`
	AgeInMonth  int            `json:"ageInMonth" db:"ageinmonth"`
	Description string         `json:"description"`
	HasMatched  bool           `json:"hasMatched"`
	ImageUrls   pq.StringArray `json:"imageUrls" db:"imageurls"`
	CreatedAt   time.Time      `json:"createdAt"`
	IsDeleted   bool           `json:"isDeleted"`
	OwnerId     string         `json:"-"`
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

type Race string

const (
	Persian          Race = "Persian"
	MaineCoon        Race = "Maine Coon"
	Siamese          Race = "Siamese"
	Ragdoll          Race = "Ragdoll"
	Bengal           Race = "Bengal"
	Sphynx           Race = "Sphynx"
	BritishShorthair Race = "British Shorthair"
	Abyssinian       Race = "Abyssinian"
	ScottishFold     Race = "Scottish Fold"
	Birman           Race = "Birman"
)
