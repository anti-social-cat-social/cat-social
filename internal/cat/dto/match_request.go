package cat

type CatMatchRequest struct {
	MatchCatId string `json:"matchCatId" binding:"required"`
	UserCatId  string `json:"userCatId" binding:"required"`
	Message    string `json:"message" binding:"required,min=5,max=120"`
}
type MatchApproveRequest struct {
	MatchId string `json:"matchId" binding:"required,uuid"`
}
