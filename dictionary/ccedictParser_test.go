package dictionary

import (
	"fmt"
	"testing"
)

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
