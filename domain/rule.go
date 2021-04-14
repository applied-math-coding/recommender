package domain

import "gorm.io/gorm"

// Rule presents: baseItemSet -> prediction
// with given confidence
type Rule struct {
	gorm.Model
	BaseItemSet    string  `gorm:"index" json:"baseItemSet"` // json of lex-ordered item-set
	Prediction     string  `json:"prediction"`               // predicted item
	Frequency      int     `json:"frequency"`                // frequency of: baseItemSet \cup {prediction}
	Confidence     float32 `json:"confidence"`               // confidence of prediction
	BaseItemLength int     `gorm:"index"`                    // length of base-item
}
