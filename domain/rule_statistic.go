package domain

type RuleStatistic = struct {
	BaseLength  int       `json:"baseLength"`
	Frequencies []int     `json:"frequencies"`
	Confidences []float32 `json:"confidences"`
}
