package tests

import (
	"main/services"
	"testing"
)

func TestFindRuleBaseLengths(t *testing.T) {
	services.InitDb(false)
	baseLength := services.FindRuleBaseLengths()
	t.Log(baseLength)
}

func TestFindStatsForBaseLength(t *testing.T) {
	services.InitDb(false)
	freqs, confs, e := services.FindStatsForBaseLength(2)
	if e != nil {
		t.Fatalf("TestFindStatsForBaseLength fails")
	}
	t.Log(freqs)
	t.Log(confs)
}

func TestComputeRuleStatistics(t *testing.T) {
	services.InitDb(false)
	ruleStats, e := services.ComputeRuleStatistics()
	if e != nil {
		t.Fatalf("TestComputeRuleStatistics fails")
	}
	t.Log(ruleStats)
}
