package domain

import "gorm.io/gorm"

type ItemCosine struct {
	gorm.Model
	LowerItem string  `gorm:"index" json:"lowerItem"`
	UpperItem string  `gorm:"index" json:"upperItem"`
	Cosine    float64 `json:"cosine"`
}
