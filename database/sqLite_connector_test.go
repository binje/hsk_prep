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

func TestGetQuestionsFromList(t *testing.T) {
	db := testSqLiteDb()
	defer db.clean()
	f1 := Fact{"Hanzi1", "Pinyin", "English"}
	f2 := Fact{"Hanzi2", "Pinyin", "English"}
	f3 := Fact{"Hanzi3", "Pinyin", "English"}
	db.InsertFact(f1)
	db.InsertFact(f2)
	db.InsertFact(f3)
	cards := db.GetQuestionsFromList([]string{})
	if len(cards) != 0 {
		t.Errorf("Expected 0 cards, got: %v", len(cards))
	}
	cards = db.GetQuestionsFromList([]string{f1.Hanzi})
	if len(cards) != 4 {
		t.Errorf("Expected 4 cards, got: %v", len(cards))
	}
	cards = db.GetQuestionsFromList([]string{f1.Hanzi, f3.Hanzi})
	if len(cards) != 8 {
		t.Errorf("Expected 8 cards, got: %v", len(cards))
	}
}
