package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/binje/hsk_prep/database"
)

type card struct {
	q string
	a string
	d string
}

func main() {

	db := database.NewSqLiteDb()
	loadVocab(db, "dictionary/hsk1vocab")

	cards := db.GetQuestions()
	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)
	i := 0
	l := len(cards)
	for len(cards) > 0 {
		r := rand.Intn(len(cards))
		c := cards[r]
		ans := strings.Split(c.Answers, ",")
		printQuestion(c)

		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		correct := false
		for _, a := range ans {
			if input == strings.TrimSpace(a) {
				i++
				fmt.Printf("Correct %v/%v\n", i, l)
				cards[r] = cards[len(cards)-1]
				cards = cards[:len(cards)-1]
				correct = true
				db.MarkKnown(c)
				break
			}
		}
		if !correct {
			fmt.Print("WRONG: ")
			fmt.Println(c.Answers)
		}
	}
}

func printQuestion(c database.Card) {
	fmt.Println()
	fmt.Println(c.Question)
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
