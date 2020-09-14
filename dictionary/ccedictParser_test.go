package dictionary

import (
	"testing"
)

func TestParse(t *testing.T) {
	facts := parseCCEDICT("cedict_1_0_ts_utf-8_mdbg.txt")
	if len(facts) != 118745 {
		t.Fail()
	}
}
