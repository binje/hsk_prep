package dictionary

import (
	"testing"
)

/*
func TestParse(t *testing.T) {
	facts := parseCCEDICT("cedict_1_0_ts_utf-8_mdbg.txt")
	if len(facts) != 118745 {
		t.Fail()
	}
}

func TestParseHSK1VocabList(t *testing.T) {
	words := ParseVocabList("hsk1")
	if len(words) != 153 {
		fmt.Println(len(words))
		fmt.Println(words)
		t.Fail()
	}
}
*/

func TestIsProperNoun(t *testing.T) {
	surname := "水 水 [Shui3] /surname Shui/"
	water := "水 水 [shui3] /water/river/liquid/beverage/additional charges or income/(of clothes) classifier for number of washes/"

	if createFact(water) == nil {
		t.Fail()
	}
	if createFact(surname) != nil {
		t.Fail()
	}
}
