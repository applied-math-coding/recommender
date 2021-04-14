package services

import (
	"main/domain"
)

func ComputeCosineStatistics() domain.CosineStatistic {
	res := domain.CosineStatistic{Cosines: make([]float64, 0)}
	DB.Raw(`select cosine
	from item_cosines`).Scan(&res.Cosines)
	return res
}
