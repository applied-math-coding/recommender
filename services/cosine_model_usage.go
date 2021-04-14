package services

import (
	"main/domain"

	"github.com/pkg/errors"
)

func MakeCosineRecommendation(itemSet domain.ItemSet) ([]domain.CosineRecommendation, error) {
	rows, e2 := DB.Raw(`
	select upper_item as item,
	cosine
	from item_cosines
	where lower_item in ?
	union
	select lower_item as item,
	cosine
	from item_cosines
	where upper_item in ?
	order by cosine desc
	limit 20`, itemSet, itemSet).Rows()
	if e2 != nil {
		return nil, errors.Wrap(e2, "DB.Raw failed")
	}
	defer rows.Close()
	res := make([]domain.CosineRecommendation, 0)
	for rows.Next() {
		recom := domain.CosineRecommendation{}
		DB.ScanRows(rows, &recom)
		res = append(res, recom)
	}
	return res, nil
}
