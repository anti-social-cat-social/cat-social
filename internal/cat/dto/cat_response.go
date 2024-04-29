package cat

import cat "1-cat-social/internal/cat/entity"

type CatUpdateResponseBody struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Race        cat.Race `json:"race"`
	Sex         cat.Sex  `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	ImageUrls   []string `json:"imageUrls"`
}
