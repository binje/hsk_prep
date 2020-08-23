package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqLiteDb struct {
	db *sql.DB
}

func NewSqLiteDb() *SqLiteDb {
	s := SqLiteDb{}
	s.init("./try.db")
	return &s
}

func (s *SqLiteDb) init(dbName string) {
	var e error
	s.db, e = sql.Open("sqlite3", dbName)
	checkError(e)
	statement, e := s.db.Prepare(`CREATE TABLE IF NOT EXISTS HSK1 
		(hanzi TEXT PRIMARY KEY, english TEXT, pinyin TEXT, 
		h2p DATETIME, h2e DATETIME, p2e DATETIME, e2h DATETIME
		)`)
	checkError(e)
	statement.Exec()
	statement.Close()
}

func (s *SqLiteDb) InsertFact(f Fact) {
	statement, e := s.db.Prepare("INSERT INTO HSK1 (english, hanzi, pinyin) VALUES (?,?,?)")
	checkError(e)
	statement.Exec(f.English, f.Hanzi, f.Pinyin)
	statement.Close()
}

func (s *SqLiteDb) InsertFacts(fs []Fact) {
	statement, e := s.db.Prepare("INSERT INTO HSK1 (english, hanzi, pinyin) VALUES (?,?,?)")
	checkError(e)
	for _, f := range fs {
		statement.Exec(f.English, f.Hanzi, f.Pinyin)
	}
	statement.Close()
}

func (s *SqLiteDb) GetQuestions() []Card {
	rows, e := s.db.Query("SELECT english, hanzi, pinyin, h2p, h2e, p2e, e2h FROM HSK1")
	checkError(e)

	cards := make([]Card, 0)
	for rows.Next() {
		var e, h, p string
		var h2p, h2e, p2e, e2h *time.Time
		if err := rows.Scan(&e, &h, &p, &h2p, &h2e, &p2e, &e2h); err != nil {
			log.Fatal(err)
		}
		if h2p == nil {
			cards = append(cards, Card{h, h, p, Hanzi, Pinyin})
		}
		if h2e == nil {
			cards = append(cards, Card{h, h, e, Hanzi, English})
		}
		if p2e == nil {
			cards = append(cards, Card{h, p, e, Pinyin, English})
		}
		if e2h == nil {
			cards = append(cards, Card{h, e, h, English, Hanzi})
		}
	}
	return cards
}

func (s *SqLiteDb) query(hanzi string) (string, string) {
	row := s.db.QueryRow("SELECT english, pinyin FROM HSK1 WHERE hanzi=?", hanzi)
	var english string
	var pinyin string
	row.Scan(&english, &pinyin)
	return english, pinyin
}

func (s *SqLiteDb) MarkKnown(c Card) {
	var field string
	switch c.QuestionType {
	case Hanzi:
		if c.AnswerType == Pinyin {
			field = "h2p"
		}
		if c.AnswerType == English {
			field = "h2e"
		}
	case Pinyin:
		if c.AnswerType == English {
			field = "p2e"
		}
	case English:
		if c.AnswerType == Hanzi {
			field = "e2h"
		}

	}
	if field == "" {
		fmt.Println("Card could not be marked known")
		fmt.Println(c)
		return
	}
	q := fmt.Sprintf(
		`UPDATE HSK1
	SET %s=datetime('now')
	WHERE hanzi="%s"`,
		field, c.Key)
	_, err := s.db.Exec(q)
	checkError(err)
}

func testSqLiteDb() SqLiteDb {
	s := SqLiteDb{}
	s.init("./tmpTest.db")
	return s
}

func (s *SqLiteDb) clean() {
	s.db.Exec("DROP TABLE HSK1")
}

func checkError(e error) {
	if e != nil {
		fmt.Println("Encounter Error: " + e.Error())
		os.Exit(1)
	}
}
