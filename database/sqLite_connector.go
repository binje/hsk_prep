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
		h2p DATETIME, h2e DATETIME, p2e DATETIME, e2h DATETIME,
		h2pPower INTEGER, h2ePower INTEGER, p2ePower INTEGER, e2hPower INTEGER
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
		now := time.Now()
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
	row := s.db.QueryRow("SELECT english, pinyin FROM HSK1 WHERE hanzi=?", hanzi)
	row.Scan(&english, &pinyin)
	return
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

	power := s.getPower(field, c.Key)

	q := fmt.Sprintf(
		`UPDATE HSK1
	SET %s=datetime('now','+%d day'), %s = %d 
	WHERE hanzi="%s"`,
		field, numDays(power), field+"Power", power+1, c.Key)
	_, err := s.db.Exec(q)
	checkError(err)
}

func (s *SqLiteDb) getPower(field, key string) int {
	var power *int
	q := fmt.Sprintf(`SELECT %s FROM HSK1 WHERE hanzi="%s"`, field+"Power", key)
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
	s.db.Exec("DROP TABLE HSK1")
}

func checkError(e error) {
	if e != nil {
		fmt.Println("Encounter Error: " + e.Error())
		os.Exit(1)
	}
}
