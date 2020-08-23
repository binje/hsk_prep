package database

import (
	"testing"
)

var _ Database = &SqLiteDb{}

func TestInsertAndGetCards(t *testing.T) {
	f := Fact{"Hanzi", "Pinyin", "English"}
	db := testSqLiteDb()
	defer db.clean()
	db.InsertFact(f)
	cards := db.GetQuestions()
	if len(cards) != 4 {
		t.Errorf("Unexpected number of cards: %v", len(cards))
	}
}

func TestKnownCardsAreNotReturned(t *testing.T) {
	f := Fact{"小姐", "Pinyin", "English"}
	db := testSqLiteDb()
	defer db.clean()
	db.InsertFact(f)
	cards := db.GetQuestions()
	if len(cards) != 4 {
		t.Errorf("Unexpected number of cards: %v", len(cards))
	}
	for i, c := range cards {
		db.MarkKnown(c)
		newCards := db.GetQuestions()
		if len(newCards) != 4-1-i {
			t.Errorf("Unexpected number of cards. want:%d, got: %d", 4-1-i, len(cards))
		}
	}
}
