package services

import (
	"main/domain"
	"math"
	"sync"
)

// Recommendations based on cosine across items.

// ItemColumn represents a sparse column
type ItemColumn = map[int]bool

func dotProduct(a ItemColumn, b ItemColumn) int {
	res := 0
	for key := range a {
		if _, ok := b[key]; ok {
			res++
		}
	}
	return res
}

func ComputeL2Norm(a ItemColumn) float64 {
	return math.Sqrt(float64(dotProduct(a, a)))
}

func CountContainingItemSets(a ItemColumn, b ItemColumn) int {
	res := len(a) + len(b)
	for row := range a {
		if _, ok := b[row]; ok {
			res = res - 1
		}
	}
	return res
}

func ExtractAndSaveCosines(itemSets []domain.ItemSet, support int, broadcaster *Broadcaster) {
	modelProgress := broadcaster.ProgressChan
	defer close(modelProgress)
	var canceled bool
	var items []string
	items, canceled = ExtractAllItems(itemSets, support, broadcaster)
	var itemColumns map[string]ItemColumn
	if !canceled {
		itemColumns, canceled = CreateItemColumns(itemSets, items, broadcaster)
	}
	var itemCosines []domain.ItemCosine
	if !canceled {
		itemCosines, canceled = ComputeCosine(items, itemColumns, broadcaster)
	}
	if canceled {
		EmitProgress("Canceled.", 1.0, domain.ProgressState.Canceled, false, modelProgress)
	} else {
		EmitProgress("Saving cosines.", 1.0, domain.ProgressState.Running, false, modelProgress)
		TruncateTable("item_cosines")
		InsertOnDb(itemCosines)
		EmitProgress("Finished.", 1.0, domain.ProgressState.Finished, false, modelProgress)
	}
}

func ComputeCosine(
	items []string,
	itemColumns map[string]ItemColumn,
	broadcaster *Broadcaster) ([]domain.ItemCosine, domain.Canceled) {
	res := make([]domain.ItemCosine, 0)
	taskIds := make([]int, 0)
	taskIdToItem := make(map[int]string)
	for taskId, item := range items {
		taskIds = append(taskIds, taskId)
		taskIdToItem[taskId] = item
	}
	computed := 0
	computedMut := sync.Mutex{}
	EmitProgress("Computing cosine: calculation", 0.0, domain.ProgressState.Running, true, broadcaster.ProgressChan)
	feedback := CreateFeedbackTag("Computing cosine: calculation", len(items), broadcaster.ProgressChan)
	parResults := DoParallel(taskIds, func(taskId int) (*ParallelResult, SkipResult) {
		lowerItem := taskIdToItem[taskId]
		if broadcaster.Cancel {
			return nil, true
		}
		internalRes := make([]domain.ItemCosine, 0)
		for _, upperItem := range items {
			if lowerItem < upperItem {
				itemColLower := itemColumns[lowerItem]
				itemColUpper := itemColumns[upperItem]
				r := dotProduct(itemColLower, itemColUpper)
				base := CountContainingItemSets(itemColLower, itemColUpper)
				cosine := float64(r) / float64(base)
				if cosine > 0.0 {
					internalRes = append(internalRes, domain.ItemCosine{LowerItem: lowerItem, UpperItem: upperItem, Cosine: cosine})
				}
			}
		}
		computedMut.Lock()
		computed++
		computedMut.Unlock()
		feedback(computed)
		return &ParallelResult{r: internalRes, taskId: taskId}, false
	})
	feedback = CreateFeedbackTag("Computing cosine: collecting results", len(parResults), broadcaster.ProgressChan)
	for idx, v := range parResults {
		feedback(idx)
		res = append(res, v.r.([]domain.ItemCosine)...)
	}
	return res, false
}

func CreateItemColumns(
	itemSets []domain.ItemSet,
	items domain.ItemSet,
	broadcaster *Broadcaster) (map[string]ItemColumn, domain.Canceled) {
	res := make(map[string]ItemColumn)
	feedback := CreateFeedbackTag("Creating item matrix: initialization", len(items), broadcaster.ProgressChan)
	for idx, item := range items {
		feedback(idx)
		if broadcaster.Cancel {
			return nil, true
		}
		res[item] = make(ItemColumn)
	}
	feedback = CreateFeedbackTag("Creating item matrix: values", len(itemSets), broadcaster.ProgressChan)
	for row, itemSet := range itemSets {
		feedback(row)
		for _, item := range itemSet {
			if broadcaster.Cancel {
				return nil, true
			}
			if _, ok := res[item]; ok {
				res[item][row] = true
			}
		}
	}
	return res, false
}

func ExtractAllItems(
	itemSets []domain.ItemSet,
	support int,
	broadcaster *Broadcaster) ([]string, domain.Canceled) {
	res := make([]string, 0)
	itemHash := make(map[string]int)
	feedback := CreateFeedbackTag("Extracting items: item hash", len(itemSets), broadcaster.ProgressChan)
	for idx, items := range itemSets {
		feedback(idx)
		for _, item := range items {
			if broadcaster.Cancel {
				return nil, true
			}
			_, ok := itemHash[item]
			if !ok {
				itemHash[item] = 0
			} else {
				itemHash[item] = itemHash[item] + 1
			}
		}
	}
	for key := range itemHash {
		if broadcaster.Cancel {
			return nil, true
		}
		if itemHash[key] >= support {
			res = append(res, key)
		}
	}
	return res, false
}
