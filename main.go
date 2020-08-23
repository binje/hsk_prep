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

	"github.com/binje/hsk_prep/utils"
)

type card struct {
	q string
	a string
	d string
}

func main() {
	// sql usage :
	/*insert("猫","cat","māo")
	insert("狗","dog","gǒu")
	english, pinyin := query("狗")
	fmt.Println(english)
	fmt.Println(pinyin)*/

	file, err := os.Open("hsk1vocab")
	if err != nil {
		fmt.Println(err)
	}
	r := csv.NewReader(file)
	r.Comma = '\t'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	cards := make([]card, len(records)*4)
	for i, r := range records {
		h := r[0]
		p := utils.TypablePinyinByLine(r[1])
		e := r[2]
		if h == "" {
			fmt.Println("h")
			fmt.Println(r)
		}
		if p == "" {
			fmt.Println("p")
			fmt.Println(r)
		}
		if e == "" {
			fmt.Println("e")
			fmt.Println(r)
		}
		j := i * 4
		cards[j] = card{h, p, "pinyin?"}
		cards[j+1] = card{h, e, "english?"}
		cards[j+2] = card{p, e, "english?"}
		cards[j+3] = card{e, h, "hanzi?"}
	}

	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)
	i := 0
	l := len(cards)
	for len(cards) > 0 {
		r := rand.Intn(len(cards))
		c := cards[r]
		ans := strings.Split(c.a, ",")
		fmt.Println()
		fmt.Println(c.q)
		fmt.Println(c.d)
		scanner.Scan()
		input := scanner.Text()
		correct := false
		for _, a := range ans {
			if input == strings.TrimSpace(a) {
				i++
				fmt.Printf("Correct %v/%v\n", i, l)
				cards[r] = cards[len(cards)-1]
				cards = cards[:len(cards)-1]
				correct = true
				break
			}
		}
		if !correct {
			fmt.Print("WRONG: ")
			fmt.Println(c.a)
		}
	}
}
