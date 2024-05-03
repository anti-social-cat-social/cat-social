package cat

import (
	cat "1-cat-social/internal/cat/entity"
	"time"

	"github.com/jinzhu/copier"
	"github.com/lib/pq"
)

type CatUpdateResponseBody struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Race        cat.Race       `json:"race"`
	Sex         cat.Sex        `json:"sex"`
	AgeInMonth  int            `json:"ageInMonth"`
	Description string         `json:"description"`
	ImageUrls   pq.StringArray `json:"imageUrls"`
}

type CatResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Race        cat.Race `json:"race"`
	Sex         cat.Sex  `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	ImageUrls   []string `json:"imageUrls"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}

func FormatCatResponse(cat *cat.Cat) CatResponse {
	catResponse := CatResponse{}
	copier.Copy(&catResponse, &cat)

	catResponse.CreatedAt = cat.CreatedAt.Format(time.RFC3339)

	return catResponse
}

func FormatCatsResponse(cats []*cat.Cat) []CatResponse {
	catsResponse := []CatResponse{}

	for _, cat := range cats {
		catResponse := FormatCatResponse(cat)
		catsResponse = append(catsResponse, catResponse)
	}

	return catsResponse
}
