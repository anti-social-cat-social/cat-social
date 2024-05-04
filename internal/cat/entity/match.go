package cat

import (
	"time"

	"github.com/lib/pq"
)

type Match struct {
	ID          string
	IssuerCatId string `db:"issuer_cat_id"`
	TargetCatId string `db:"target_cat_id"`
	Message     string
	Status      MatchStatus
	CreatedAt   time.Time
	IsDeleted   bool
	IssuedBy    string
}

type MatchStatus string

const (
	Submitted MatchStatus = "submitted"
	Cancelled MatchStatus = "cancelled"
	Approved  MatchStatus = "approved"
	Rejected  MatchStatus = "rejected"
)

type MatchDetail struct {
	MatchesID string `json:"matches_id" db:"matches_id"`
	UsersName string `json:"users_name" db:"users_name"`
	UsersEmail string `json:"users_email" db:"users_email"`
	IssuerCatID string `json:"issuer_cat_id" db:"issuer_cat_id"`
	CatsIssuerName string `json:"cats_issuer_name" db:"cats_issuer_name"`
	CatsIssuerRace string `json:"cats_issuer_race" db:"cats_issuer_race"`
	CatsIssuerSex string `json:"cats_issuer_sex" db:"cats_issuer_sex"`
	CatsIssuerDescription string `json:"cats_issuer_description" db:"cats_issuer_description"`
	CatsIssuerAgeinmonth int `json:"cats_issuer_ageinmonth" db:"cats_issuer_ageinmonth"`
	CatsIssuerImageurls pq.StringArray `json:"cats_issuer_imageurls" db:"cats_issuer_imageurls"`
	CatsIssuerHasmatched bool `json:"cats_issuer_hasmatched" db:"cats_issuer_hasmatched"`
	CatsIssuerCreatedat time.Time `json:"cats_issuer_createdat" db:"cats_issuer_createdat"`
	MatchesMessage string `json:"matches_message" db:"matches_message"`
	MatchesCreatedat time.Time `json:"matches_createdat" db:"matches_createdat"`
	TargetCatID string `json:"target_cat_id" db:"target_cat_id"`
	CatsTargetName string `json:"cats_target_name" db:"cats_target_name"`
	CatsTargetRace string `json:"cats_target_race" db:"cats_target_race"`
	CatsTargetSex string `json:"cats_target_sex" db:"cats_target_sex"`
	CatsTargetDescription string `json:"cats_target_description" db:"cats_target_description"`
	CatsTargetAgeinmonth int `json:"cats_target_ageinmonth" db:"cats_target_ageinmonth"`
	CatsTargetImageurls pq.StringArray `json:"cats_target_imageurls" db:"cats_target_imageurls"`
	CatsTargetHasmatched bool `json:"cats_target_hasmatched" db:"cats_target_hasmatched"`
	CatsTargetCreatedat time.Time `json:"cats_target_createdat" db:"cats_target_createdat"`
}