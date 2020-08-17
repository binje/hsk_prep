package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"
)

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

	m := make(map[string]string)
	for _, r := range records {
		m[r[0]] = r[2]
	}

	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)
	for {
		keys := reflect.ValueOf(m).MapKeys()
		q := keys[rand.Intn(len(keys))].String()
		a := m[q]
		fmt.Print(q)
		scanner.Scan()
		input := scanner.Text()
		if input == a {
			fmt.Println("Correct")
		} else {
			fmt.Println("WRONG: ")
			fmt.Println(input)
			fmt.Println(a)
		}
	}
}
