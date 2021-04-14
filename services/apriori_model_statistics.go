package services

import (
	"main/domain"

	"github.com/pkg/errors"
)

func ComputeRuleStatistics() ([]domain.RuleStatistic, error) {
	res := make([]domain.RuleStatistic, 0)
	for _, bl := range FindRuleBaseLengths() {
		stat := domain.RuleStatistic{BaseLength: bl}
		freqs, confs, e := FindStatsForBaseLength(bl)
		if e != nil {
			return nil, errors.WithStack(e)
		}
		stat.Frequencies = freqs
		stat.Confidences = confs
		res = append(res, stat)
	}
	return res, nil
}

func FindRuleBaseLengths() []int {
	res := make([]int, 0)
	DB.Raw(`select base_item_length
					from rules
					group by base_item_length`).Scan(&res)
	return res
}

func FindStatsForBaseLength(bl int) ([]int, []float32, error) {
	rows, e := DB.Raw(`select
	frequency,
	confidence
	from rules
	where base_item_length = ?`, bl).Rows()
	if e != nil {
		return nil, nil, errors.Wrap(e, "DB.Raw failed")
	}
	defer rows.Close()
	frequencies := make([]int, 0)
	confidences := make([]float32, 0)
	for rows.Next() {
		var freq int
		var conf float32
		rows.Scan(&freq, &conf)
		frequencies = append(frequencies, freq)
		confidences = append(confidences, conf)
	}
	return frequencies, confidences, nil
}
