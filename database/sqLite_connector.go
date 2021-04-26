package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqLiteDb struct {
	db *sql.DB
}

func NewSqLiteDb() *SqLiteDb {
	s := SqLiteDb{}
	s.init("./dict.db")
	return &s
}

func (s *SqLiteDb) init(dbName string) {
	var e error
	s.db, e = sql.Open("sqlite3", dbName)
	checkError(e)
	statement, e := s.db.Prepare(`CREATE TABLE IF NOT EXISTS CHINESE 
		(hanzi TEXT PRIMARY KEY, english TEXT, pinyin TEXT, 
		h2p DATETIME, h2e DATETIME, p2e DATETIME, e2h DATETIME,
		h2pPower INTEGER, h2ePower INTEGER, p2ePower INTEGER, e2hPower INTEGER
		)`)
	checkError(e)
	statement.Exec()
	statement.Close()
}

func (s *SqLiteDb) InsertFact(f Fact) {
	statement, e := s.db.Prepare("INSERT INTO CHINESE (english, hanzi, pinyin) VALUES (?,?,?)")
	checkError(e)
	statement.Exec(f.English, f.Hanzi, f.Pinyin)
	statement.Close()
}

func (s *SqLiteDb) InsertFacts(fs []Fact) {
	tx, err := s.db.Begin()
	if err != nil {
		panic(err)
	}
	statement, e := tx.Prepare("INSERT INTO CHINESE (english, hanzi, pinyin) VALUES (?,?,?)")
	checkError(e)
	lastKey := ""
	for _, f := range fs {
		//Uniqueness needed. TODO glob together
		if f.Hanzi == lastKey {
			continue
		}
		lastKey = f.Hanzi
		_, err = statement.Exec(f.English, f.Hanzi, f.Pinyin)
		if err != nil {
			panic(err)
		}
	}
	err = statement.Close()
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (s *SqLiteDb) GetQuestions() []Card {
	rows, e := s.db.Query("SELECT english, hanzi, pinyin, h2p, h2e, p2e, e2h FROM CHINESE")
	checkError(e)
	return makeCards(rows)

}

func (s *SqLiteDb) GetQuestionsFromList(words []string) []Card {
	if len(words) == 0 {
		return []Card{}
	}

	q := fmt.Sprintf(
		`SELECT english, hanzi, pinyin, h2p, h2e, p2e, e2h FROM CHINESE WHERE hanzi in ("%s")`, strings.Join(words, "\",\""))
	rows, e := s.db.Query(q)
	checkError(e)
	return makeCards(rows)
}

func makeCards(rows *sql.Rows) []Card {
	cards := make([]Card, 0)
	now := time.Now()
	for rows.Next() {
		var e, h, p string
		var h2p, h2e, p2e, e2h *time.Time
		if err := rows.Scan(&e, &h, &p, &h2p, &h2e, &p2e, &e2h); err != nil {
			log.Fatal(err)
		}
		if h2p == nil || now.After(*h2p) {
			cards = append(cards, Card{h, h, p, Hanzi, Pinyin})
		}
		if h2e == nil || now.After(*h2e) {
			cards = append(cards, Card{h, h, e, Hanzi, English})
		}
		if p2e == nil || now.After(*p2e) {
			cards = append(cards, Card{h, p, e, Pinyin, English})
		}
		if e2h == nil || now.After(*e2h) {
			cards = append(cards, Card{h, e, h, English, Hanzi})
		}
	}
	return cards

}

func (s *SqLiteDb) query(hanzi string) (english, pinyin string) {
	row := s.db.QueryRow("SELECT english, pinyin FROM CHINESE WHERE hanzi=?", hanzi)
	row.Scan(&english, &pinyin)
	return
}

func (s *SqLiteDb) MarkKnown(c Card) {

	field := getFieldName(c)
	if field == "" {
		fmt.Println("Card could not be marked known")
		fmt.Println(c)
		return
	}

	power := s.getPower(field, c.Key)

	q := fmt.Sprintf(
		`UPDATE CHINESE
	SET %s=datetime('now','+%d day'), %s = %d 
	WHERE hanzi="%s"`,
		field, numDays(power), field+"Power", power+1, c.Key)
	_, err := s.db.Exec(q)
	checkError(err)
}

func (s *SqLiteDb) MarkUnknown(c Card) {
	field := getFieldName(c)
	if field == "" {
		fmt.Println("Card could not be marked unknown")
		fmt.Println(c)
		return
	}

	q := fmt.Sprintf(
		`UPDATE CHINESE
	SET %s=datetime('now'), %s = %d 
	WHERE hanzi="%s"`,
		field, field+"Power", 0, c.Key)
	_, err := s.db.Exec(q)
	checkError(err)
}

func getFieldName(c Card) string {
	switch c.QuestionType {
	case Hanzi:
		if c.AnswerType == Pinyin {
			return "h2p"
		}
		if c.AnswerType == English {
			return "h2e"
		}
	case Pinyin:
		if c.AnswerType == English {
			return "p2e"
		}
	case English:
		if c.AnswerType == Hanzi {
			return "e2h"
		}

	}
	return ""
}

func (s *SqLiteDb) getPower(field, key string) int {
	var power *int
	q := fmt.Sprintf(`SELECT %s FROM CHINESE WHERE hanzi="%s"`, field+"Power", key)
	row := s.db.QueryRow(q)
	row.Scan(&power)
	if power == nil {
		return 0
	}
	return *power
}

func numDays(n int) int {
	d := 1
	for i := 0; i < n; i++ {
		d *= 5
	}
	return d
}

func testSqLiteDb() SqLiteDb {
	s := SqLiteDb{}
	s.init("./tmpTest.db")
	return s
}

func (s *SqLiteDb) clean() {
	s.db.Exec("DROP TABLE CHINESE")
}

func checkError(e error) {
	if e != nil {
		fmt.Println("Encounter Error: " + e.Error())
		os.Exit(1)
	}
}
