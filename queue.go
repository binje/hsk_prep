package main

import (
	"math/rand"
	"time"

	"github.com/binje/hsk_prep/database"
)

type Queue struct {
	q     [][]database.Card
	cache []database.Card
	i     int
}

func GenerateQueue(cards []database.Card) *Queue {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	q := make([][]database.Card, 3)
	for _, card := range cards {
		q[card.AnswerType] = append(q[card.AnswerType], card)
	}
	return &Queue{
		q:     q,
		cache: make([]database.Card, 0),
	}
}

func (q *Queue) GetNext() database.Card {
	if len(q.cache) == 0 {
		q.i++
		q.i %= 3
		q.cache = q.q[q.i][:5]
		q.q[q.i] = q.q[q.i][6:]
	}
	c := q.cache[0]
	q.cache = q.cache[1:]
	return c
}

func (q *Queue) HasNext() bool {
	if len(q.cache) != 0 {
		return true
	}
	for _, cards := range q.q {
		if len(cards) != 0 {
			return true
		}
	}
	return false
}

func (q *Queue) MarkUnknown(c database.Card) {
	at := c.AnswerType
	q.q[at] = append(q.q[at], c)
	q.q[at][0], q.q[at][len(q.q[at])-1] = q.q[at][len(q.q[at])-1], q.q[at][0]
}
