package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/binje/hsk_prep/database"
	"github.com/binje/hsk_prep/dictionary"
)

func main() {

	db := database.NewSqLiteDb()
	loadCCEDICT(db)

	cards := db.GetQuestionsFromList(dictionary.ParseVocabList("dictionary/hsk1"))
	scanner := bufio.NewScanner(os.Stdin)
	i := 0
	l := len(cards)
	q := GenerateQueue(cards)

	for q.HasNext() {
		c := q.GetNext()
		printQuestion(c)

		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if isCorrect(input, c.Answers) {
			i++
			fmt.Printf("Correct %v/%v\n", i, l)
			db.MarkKnown(c)
		} else {
			fmt.Print("WRONG: ")
			fmt.Println(c.Answers)
			db.MarkUnknown(c)
			q.MarkUnknown(c)
		}
	}
	fmt.Println("CONGRATULATIONS! YOU KNOW ALL OF THE WORDS!")
}

func isCorrect(input, answers string) bool {
	if input == "" {
		return false
	}
	input = strings.ToLower(input)
	ans := strings.Split(answers, "/")
	for _, a := range ans {
		if input == strings.ToLower(strings.TrimSpace(a)) {
			return true
		}
	}
	return false
}

func printQuestion(c database.Card) {
	fmt.Println()
	switch c.QuestionType {
	case database.Hanzi:
		fmt.Println(c.Question)
	case database.Pinyin:
		fmt.Println(c.Question)
	case database.English:
		definitions := strings.Split(c.Question, "/")
		for _, d := range definitions {
			if d != "" {
				fmt.Println(d)
			}
		}
	}
	switch c.AnswerType {
	case database.Hanzi:
		fmt.Println("Hanzi?")
	case database.Pinyin:
		fmt.Println("Pinyin?")
	case database.English:
		fmt.Println("English?")
	}
}

func loadVocab(db database.Database, filePath string) {
	file, err := os.Open("dictionary/hsk1vocab")
	if err != nil {
		fmt.Println(err)
	}
	r := csv.NewReader(file)
	r.Comma = '\t'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range records {
		db.InsertFact(database.Fact{r[0], r[1], r[2]})
	}
}

func loadCCEDICT(db database.Database) {
	facts := dictionary.ParseCCEDICT()
	db.InsertFacts(facts)
}
