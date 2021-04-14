package domain

type Recommendation = struct {
	Prediction string  `json:"prediction"` // predicted item
	Frequency  int     `json:"frequency"`  // frequency of: baseItemSet \cup {prediction}
	Confidence float32 `json:"confidence"` // confidence of prediction
}
