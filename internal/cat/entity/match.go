package cat

import (
	"time"
)

type Match struct {
	ID            string
	IssuerCatId   string `db:"issuer_cat_id"`
	TargetCatId   string `db:"target_cat_id"`
	Message       string
	Status        MatchStatus
	CreatedAt     time.Time
	IsDeleted     bool
	IssuedBy      string
	TargetOwnerID string `db:"target_cat_owner"`
}

type MatchStatus string

const (
	Submitted MatchStatus = "submitted"
	Cancelled MatchStatus = "cancelled"
	Approved  MatchStatus = "approved"
	Rejected  MatchStatus = "rejected"
)
