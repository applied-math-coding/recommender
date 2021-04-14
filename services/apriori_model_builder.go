package services

import (
	"encoding/json"
	"fmt"
	"main/domain"
	"main/utils"
	"sort"
	"strconv"

	"github.com/pkg/errors"
)

func ExtractAndSaveRules(
	itemSets []domain.ItemSet,
	support int,
	broadcaster *Broadcaster) {
	modelProgress := broadcaster.ProgressChan
	defer close(modelProgress)
	rules, canceled, e := ExtractRules(
		itemSets,
		support,
		nil,
		nil,
		broadcaster)
	if e != nil {
		HandleError(errors.WithStack(e), nil, false)
		EmitProgress("Error when extracting rules.", 0.0, domain.ProgressState.Error, false, modelProgress)
		return
	}
	if canceled {
		EmitProgress("Canceled.", 1.0, domain.ProgressState.Canceled, false, modelProgress)
	} else {
		EmitProgress("Saving rules.", 1.0, domain.ProgressState.Running, false, modelProgress)
		TruncateTable("rules")
		InsertOnDb(rules)
		EmitProgress("Finished.", 1.0, domain.ProgressState.Finished, false, modelProgress)
	}
}

// ExtractRules returns association-rules based on the given itemSets.
// The algorithm is heavily based on storing the items in lex-order.
// Iteration is made along length of item-sets and only considers those item-sets
// which fulfill required support
func ExtractRules(
	itemSets []domain.ItemSet,
	support int,
	singleItemFrequencies []domain.ItemFrequency,
	previousItemFrequencies []domain.ItemFrequency,
	broadcaster *Broadcaster) ([]domain.Rule, domain.Canceled, error) {
	frequencies := make([]domain.ItemFrequency, 0)
	if singleItemFrequencies == nil { // at start
		var canceled bool
		singleItemFrequencies, canceled = ExtractSingleItemFrequencies(itemSets, broadcaster)
		if canceled {
			return nil, true, nil
		}
		singleItemFrequencies, canceled = FilterItemFrequencies(singleItemFrequencies, support, broadcaster)
		if canceled {
			return nil, true, nil
		}
		frequencies = singleItemFrequencies
	} else {
		feedback := CreateFeedbackTag(
			fmt.Sprintf("Creating next level item-sets (base-items: %d)", len(previousItemFrequencies)),
			len(previousItemFrequencies),
			broadcaster.ProgressChan)
		for idx, baseFreq := range previousItemFrequencies { // search further frequent items
			if broadcaster.Cancel {
				return nil, true, nil
			}
			feedback(idx)
			for _, itemSet := range CreateItemSets(baseFreq.ItemSet, singleItemFrequencies, broadcaster) {
				if broadcaster.Cancel {
					return nil, true, nil
				}
				itemFreq, canceled := CreateItemFrequency(itemSets, &baseFreq, itemSet, broadcaster)
				if canceled {
					return nil, true, nil
				}
				if itemFreq.Frequency >= support {
					frequencies = append(frequencies, *itemFreq)
				}
			}
		}
	}
	rules := make([]domain.Rule, 0)
	if len(frequencies) > 0 { // create rules and search in next level (item-set length)
		var e error
		var canceled bool
		previousFreqHash := make(map[string]int)
		for _, freq := range previousItemFrequencies {
			if canceled {
				return nil, true, nil
			}
			itemSetJson, _ := json.Marshal(freq.ItemSet)
			previousFreqHash[string(itemSetJson)] = freq.Frequency
		}
		rules, canceled, e = createRules(frequencies, previousFreqHash, broadcaster)
		if canceled {
			return nil, true, nil
		}
		if e != nil {
			return nil, false, errors.WithStack(e)
		}
		furtherRules, canceled, e := ExtractRules(
			itemSets,
			support,
			singleItemFrequencies,
			frequencies,
			broadcaster)
		if canceled {
			return nil, true, nil
		}
		if e != nil {
			return nil, false, errors.WithStack(e)
		}
		rules = append(rules, furtherRules...)
	}
	return rules, false, nil
}

func createRules(
	freqs []domain.ItemFrequency,
	prevFreqHash map[string]int,
	broadcaster *Broadcaster) ([]domain.Rule, domain.Canceled, error) {
	res := make([]domain.Rule, 0)
	feedback := CreateFeedbackTag(
		"Creating rules for item-sets of length "+strconv.Itoa(len(freqs[0].ItemSet))+".",
		len(freqs),
		broadcaster.ProgressChan)
	for idx, freq := range freqs {
		if broadcaster.Cancel {
			return make([]domain.Rule, 0), true, nil
		}
		feedback(idx)
		if len(freq.ItemSet) > 1 {
			for itemIdx, item := range freq.ItemSet {
				if broadcaster.Cancel {
					return make([]domain.Rule, 0), true, nil
				}
				baseItemSet := make(domain.ItemSet, len(freq.ItemSet)-1)
				copy(baseItemSet, freq.ItemSet[:itemIdx])
				copy(baseItemSet[itemIdx:], freq.ItemSet[itemIdx+1:])
				sort.Strings(baseItemSet)
				baseItemSetJson, jsonError := json.Marshal(baseItemSet)
				if jsonError != nil {
					return nil, false, errors.Wrap(jsonError, "json.Marshal failed")
				}
				baseFrequency, ok := prevFreqHash[string(baseItemSetJson)]
				if ok {
					confidence := float32(freq.Frequency) / float32(baseFrequency)
					if confidence > 0.5 {
						rule := domain.Rule{
							BaseItemSet:    string(baseItemSetJson),
							Prediction:     item,
							Frequency:      freq.Frequency,
							Confidence:     confidence,
							BaseItemLength: len(freq.ItemSet) - 1}
						res = append(res, rule)
					}
				}
			}
		} else {
			baseItemSetJson, jsonError := json.Marshal(freq.ItemSet[:0])
			if jsonError != nil {
				return nil, false, errors.Wrap(jsonError, "json.Marshal failed")
			}
			rule := domain.Rule{
				BaseItemSet:    string(baseItemSetJson),
				Prediction:     freq.ItemSet[0],
				Frequency:      freq.Frequency,
				Confidence:     1.0,
				BaseItemLength: 1}
			res = append(res, rule)
		}
	}
	return res, false, nil
}

// CreateItemSets creates a list of ItemSet which are constructed from 'base' by added
// the itemSet one after each other from singleItemFrequencies. Only those itemSet
// from singleItemFrequencies are considered which are lex-higher than the highest
// in 'base'. 'base' is expected to be lex-ordered.
func CreateItemSets(
	base domain.ItemSet,
	singleItemFrequencies []domain.ItemFrequency,
	broadcaster *Broadcaster) []domain.ItemSet {
	res := make([]domain.ItemSet, 0)
	highest := base[len(base)-1]
	for _, freq := range singleItemFrequencies {
		if broadcaster.Cancel {
			return make([]domain.ItemSet, 0)
		}
		item := freq.ItemSet[0]
		if item > highest {
			newItem := make(domain.ItemSet, len(base))
			copy(newItem, base)
			newItem = append(newItem, item)
			res = append(res, newItem)
		}
	}
	return res
}

// FilterFrequentItemSets filters on those items which have at least frequency of 'support'
// and return those itemFrequencies which meet the support and those which do not.
func FilterItemFrequencies(
	itemFreqs []domain.ItemFrequency,
	support int,
	broadcaster *Broadcaster) ([]domain.ItemFrequency, domain.Canceled) {
	withSupport := make([]domain.ItemFrequency, 0)
	feedback := CreateFeedbackTag("Searching single frequent items", len(itemFreqs), broadcaster.ProgressChan)
	for idx, item := range itemFreqs {
		if broadcaster.Cancel {
			return nil, true
		}
		feedback(idx)
		if item.Frequency >= support {
			withSupport = append(withSupport, item)
		}
	}
	return withSupport, false
}

// ExtractSingleItemFrequencies extract all different items and returns
// ItemFrequencies for them. The returned list is lex-sorted.
func ExtractSingleItemFrequencies(
	itemSets []domain.ItemSet,
	broadcaster *Broadcaster) ([]domain.ItemFrequency, domain.Canceled) {
	itemMap := make(map[string]*domain.ItemFrequency)
	feedback := CreateFeedbackTag("Searching single frequent items", len(itemSets), broadcaster.ProgressChan)
	for idx, itemSet := range itemSets {
		feedback(idx)
		for _, item := range itemSet {
			if broadcaster.Cancel {
				return nil, true
			}
			_, ok := itemMap[item]
			if !ok {
				itemMap[item] = &domain.ItemFrequency{
					ItemSet:   domain.ItemSet{item},
					Frequency: 0}
			}
			itemFreq := itemMap[item]
			itemFreq.Frequency = itemFreq.Frequency + 1
		}
	}
	items := make([]string, 0)
	for item := range itemMap {
		items = append(items, item)
	}
	sort.Strings(items)
	res := make([]domain.ItemFrequency, 0)
	for _, item := range items {
		res = append(res, *itemMap[item])
	}
	return res, false
}

// CreateItemFrequency creates an ItemFrequency for given ItemSet based on
// counting in given ItemSets. The baseItemFreq is the one based on which the given
// itemSet is created.
func CreateItemFrequency(
	itemSets []domain.ItemSet,
	baseItemFreq *domain.ItemFrequency,
	itemSet domain.ItemSet,
	broadcaster *Broadcaster) (*domain.ItemFrequency, domain.Canceled) {
	res := domain.ItemFrequency{
		ItemSet:   itemSet,
		Frequency: 0}
	for _, testItemSet := range itemSets {
		if broadcaster.Cancel {
			return nil, true
		}
		if utils.IsContained(itemSet, testItemSet) {
			res.Frequency = res.Frequency + 1
		}
	}
	return &res, false
}
