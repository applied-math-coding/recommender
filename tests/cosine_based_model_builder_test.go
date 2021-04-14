package tests

import (
	"main/domain"
	"main/services"
	"testing"
)

var cosineItemSets = []domain.ItemSet{
	{"A", "B", "C", "D"},
	{"A", "B"},
	{"B", "C", "D"},
	{"A", "B", "D"},
	{"B", "C"},
	{"C", "D"},
	{"B", "D"}}

func TestExtractAllItems(t *testing.T) {
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	items, _ := services.ExtractAllItems(cosineItemSets, 0, broadcaster)
	t.Log(items)
	if len(items) != 4 {
		t.Fatalf("TestExtractAllItems fails")
	}
}

func TestCreateItemMatrix(t *testing.T) {
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	items, _ := services.ExtractAllItems(cosineItemSets, 0, broadcaster)
	t.Log(items)
	matrix, _ := services.CreateItemColumns(cosineItemSets, items, broadcaster)
	t.Log(matrix)
}

func TestComputeCosine(t *testing.T) {
	broadcaster := services.CreateBroadcaster(make(chan domain.ProgressMessage))
	items, _ := services.ExtractAllItems(cosineItemSets, 0, broadcaster)
	itemColumns, _ := services.CreateItemColumns(cosineItemSets, items, broadcaster)
	cosines, _ := services.ComputeCosine(items, itemColumns, broadcaster)
	for _, c := range cosines {
		t.Logf(`%v, %v: %v`, c.LowerItem, c.UpperItem, c.Cosine)
	}
}
