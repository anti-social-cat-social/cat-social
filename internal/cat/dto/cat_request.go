package cat

import (
	cat "1-cat-social/internal/cat/entity"
)

// Oneof
// - "Persian"
// - "Maine Coon"
// - "Siamese"
// - "Ragdoll"
// - "Bengal"
// - "Sphynx"
// - "British Shorthair"
// - "Abyssinian"
// - "Scottish Fold"
// - "Birman" */

// - "Persian"
// - "Siamese"
// - "Ragdoll"
// - "Bengal"
// - "Sphynx"
// - "Abyssinian"
// - "Birman" */
// - "Scottish Fold"
// - "Maine Coon"
// - "British Shorthair"
type CatUpdateRequestBody struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        cat.Race `json:"race" validate:"required,oneof=Persian Siamese Ragdoll Bengal Sphynx Abyssinian Birman 'Scottish Fold' 'Maine Coon' 'British Shorthair'"`
	Sex         cat.Sex  `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,dive,required"`
}
