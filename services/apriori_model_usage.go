package services

import (
	"encoding/json"
	"main/domain"
	"sort"

	"github.com/pkg/errors"
)

func MakeAprioriRecommendation(itemSet domain.ItemSet) ([]domain.Recommendation, error) {
	sort.Strings(itemSet)
	baseItemSet, e1 := json.Marshal(itemSet)
	if e1 != nil {
		return nil, errors.Wrap(e1, "json.Marshal failed")
	}
	rows, e2 := DB.Raw(`select
	prediction,
	frequency,
	confidence
	from rules
	where base_item_set = ?
	order by frequency desc, confidence desc
	limit 10`, baseItemSet).Rows()
	if e2 != nil {
		return nil, errors.Wrap(e2, "DB.Raw failed")
	}
	defer rows.Close()
	res := make([]domain.Recommendation, 0)
	for rows.Next() {
		recom := domain.Recommendation{}
		DB.ScanRows(rows, &recom)
		res = append(res, recom)
	}
	return res, nil
}

func FindExampleRules() ([]domain.Rule, error) {
	baseLenghts, e1 := FindAllRuleLengths()
	if e1 != nil {
		return nil, errors.WithStack(e1)
	}
	res := make([]domain.Rule, 0)
	for _, baseLength := range baseLenghts {
		rules := FindFirstNRules(baseLength, 3)
		res = append(res, rules...)
	}
	return res, nil
}

func FindFirstNRules(length int, n int) []domain.Rule {
	res := make([]domain.Rule, 0)
	DB.Limit(n).Where("base_item_length = ?", length).Find(&res)
	return res
}

func FindAllRuleLengths() ([]int, error) {
	rows, e1 := DB.Raw(`select
	base_item_length
	from rules
	group by base_item_length`).Rows()
	if e1 != nil {
		return nil, errors.Wrap(e1, "DB.Raw failed")
	}
	defer rows.Close()
	res := make([]int, 0)
	for rows.Next() {
		var length int
		rows.Scan(&length)
		res = append(res, length)
	}
	return res, nil
}
