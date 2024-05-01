package cat

import (
	cat "1-cat-social/internal/cat/entity"

	"github.com/lib/pq"
)

type CatUpdateRequestBody struct {
	Name        string         `json:"name" validate:"required,min=1,max=30"`
	Race        cat.Race       `json:"race" validate:"required,oneof=Persian Siamese Ragdoll Bengal Sphynx Abyssinian Birman 'Scottish Fold' 'Maine Coon' 'British Shorthair'"`
	Sex         cat.Sex        `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int            `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string         `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   pq.StringArray `json:"imageUrls" validate:"required,gt=0,dive,url,required"`
}
