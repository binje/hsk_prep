package utils

import "testing"

func TestTypeablePinyinWords(t *testing.T) {
	translation := []struct {
		in  string
		out string
	}{
		{"dú", "du2"},
		{"méi", "mei2"},
		{"shǎo", "shao3"},
		{"zì", "zi4"},
		{"nǚ", "nü3"},
		{"zhuō", "zhuo1"},
	}
	for _, tr := range translation {
		if got := TypablePinyin(tr.in); got != tr.out {
			t.Errorf("Incorrect pinyin translation, got: %s, want: %s", got, tr.out)
		}
	}
}

func TestTypable2Pinyin(t *testing.T) {
	translation := []struct {
		out string
		in  string
	}{
		{"dú", "du2"},
		{"méi", "mei2"},
		{"shǎo", "shao3"},
		{"zì", "zi4"},
		{"nǚ", "nü3"},
		{"zhuō", "zhuo1"},
	}
	for _, tr := range translation {
		if got := Typable2Pinyin(tr.in); got != tr.out {
			t.Errorf("Incorrect pinyin translation, got: %s, want: %s", got, tr.out)
		}
	}
}
