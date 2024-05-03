package cat

import (
	cat "1-cat-social/internal/cat/entity"

	"github.com/lib/pq"
)

type CatUpdateRequestBody struct {
	Name        string         `json:"name" validate:"required,min=1,max=30,valid_name"`
	Race        cat.Race       `json:"race" validate:"required,oneof=Persian Siamese Ragdoll Bengal Sphynx Abyssinian Birman 'Scottish Fold' 'Maine Coon' 'British Shorthair'"`
	Sex         cat.Sex        `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int            `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string         `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   pq.StringArray `json:"imageUrls" validate:"required,gt=0,dive,url,required"`
}

type CatRequestQueryParams struct {
	ID         string `form:"id"`
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
	Race       string `form:"race" binding:"omitempty,oneof='Persian' 'MaineCoon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex        string `form:"sex" binding:"omitempty,oneof=male female"`
	HasMatched string `form:"hasMatched" binding:"omitempty,oneof=true false"`
	AgeInMonth string `form:"ageInMonth"`
	Owned      string `form:"owned" binding:"omitempty,oneof=true false"`
	Search     string `form:"search"`
}
