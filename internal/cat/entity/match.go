package cat

import (
	"time"
)

type Match struct {
	ID          string
	IssuerCatId string
	TargetCatId string
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
