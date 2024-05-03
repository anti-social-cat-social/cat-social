package cat

import entity "1-cat-social/internal/cat/entity"

type MatchResponse struct {
	ID             string     `json:"id" db:"id"`
	IssueBy        issuedBy   `json:"issuedBy" db:"issuedBy"`
	MatchCatDetail entity.Cat `json:"matchCatDetail" db:"matchCatDetail"`
	UserCatDetail  entity.Cat `json:"userCatDetail" db:"userCatDetail"`
	Message        string     `json:"message" db:"message"`
	CreatedAt      string     `json:"createdAt" db:"createdat"`
}

type issuedBy struct {
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	CreatedAt string `json:"createdAt" db:"createdat"`
}
