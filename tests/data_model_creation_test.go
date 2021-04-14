package tests

import (
	"main/domain"
	"main/services"
	"reflect"
	"testing"
)

var itemSets = []domain.ItemSet{
	{"A", "B", "C", "D"},
	{"A", "B"},
	{"B", "C", "D"},
	{"A", "B", "D"},
	{"B", "C"},
	{"C", "D"},
	{"B", "D"}}

var singleItemFrequencies = []domain.ItemFrequency{
	{ItemSet: domain.ItemSet{"A"}, Frequency: 3},
	{ItemSet: domain.ItemSet{"B"}, Frequency: 6},
	{ItemSet: domain.ItemSet{"C"}, Frequency: 4},
	{ItemSet: domain.ItemSet{"D"}, Frequency: 5}}

func TestExtractSingleItemFrequencies(t *testing.T) {
	channel := make(chan domain.ProgressMessage)
	defer close(channel)
	go func() {
		for range channel {
		}
	}()
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	freqs, _ := services.ExtractSingleItemFrequencies(itemSets, broadcaster)
	for idx := range freqs {
		if !reflect.DeepEqual(&freqs[idx], &singleItemFrequencies[idx]) {
			t.Fatalf("TestExtractSingleItemFrequencies fails")
		}
	}
}

func TestFilterFrequentItemSets(t *testing.T) {
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	freqs, _ := services.FilterItemFrequencies(singleItemFrequencies, 4, broadcaster)
	expected := []domain.ItemFrequency{
		{ItemSet: domain.ItemSet{"B"}, Frequency: 6},
		{ItemSet: domain.ItemSet{"C"}, Frequency: 4},
		{ItemSet: domain.ItemSet{"D"}, Frequency: 5}}
	for idx := range freqs {
		if !reflect.DeepEqual(&freqs[idx], &expected[idx]) {
			t.Fatalf("TestFilterFrequentItemSets fails")
		}
	}
}

func TestCreateItemSets(t *testing.T) {
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	base := domain.ItemSet{"A", "B"}
	itemSets := services.CreateItemSets(base, singleItemFrequencies, broadcaster)
	expected := []domain.ItemSet{
		{"A", "B", "C"},
		{"A", "B", "D"}}
	if !reflect.DeepEqual(itemSets, expected) {
		t.Fatalf("TestCreateItemSets fails")
	}
}

func TestCreateItemFrequency(t *testing.T) {
	baseFreq := domain.ItemFrequency{
		ItemSet:   domain.ItemSet{"B"},
		Frequency: 6}
	itemSet := domain.ItemSet{"B", "C"}
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	freq, _ := services.CreateItemFrequency(itemSets, &baseFreq, itemSet, broadcaster)
	expected := domain.ItemFrequency{
		ItemSet:   domain.ItemSet{"B", "C"},
		Frequency: 3}
	if !reflect.DeepEqual(*freq, expected) {
		t.Fatalf("TestCreateItemFrequency fails")
	}
}

/**
A,B 3
B,C 3
B,D 4
C,D 3

[]->A (rule)
[]->B (rule)
[]->C (rule)
[]->D (rule)
A->B 3/3 (rule)
B->A 3/6
B->C 3/6
C->B 3/4 (rule)
B->D 4/6 (rule)
D->B 4/5 (rule)
C->D 3/3 (rule)
D->C 3/5 (rule)
*/
func TestExtractRules(t *testing.T) {
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	rules, _, e := services.ExtractRules(itemSets, 3, nil, nil, broadcaster)
	if e != nil || len(rules) != 10 {
		t.Fatalf("TestExtractRules fails")
	}
}
