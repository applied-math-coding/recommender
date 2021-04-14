package domain

type CosineRecommendation = struct {
	Item   string  `json:"item"`
	Cosine float64 `json:"cosine"`
}
