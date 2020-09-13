package database

type Fact struct {
	Hanzi   string
	Pinyin  string
	English string
}

type FactType int

const (
	Hanzi FactType = iota
	Pinyin
	English
)

type Card struct {
	Key          string
	Question     string
	Answers      string
	QuestionType FactType
	AnswerType   FactType
}

type Database interface {
	InsertFact(Fact)
	InsertFacts([]Fact)
	GetQuestions() []Card
	GetQuestionsFromList([]string) []Card
	MarkKnown(Card)
}
