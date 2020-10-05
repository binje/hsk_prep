package dictionary

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/binje/hsk_prep/database"
)

func ParseCCEDICT() []database.Fact {
	return parseCCEDICT("dictionary/cedict_1_0_ts_utf-8_mdbg.txt")
}

func parseCCEDICT(filePath string) []database.Fact {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	facts := make([]database.Fact, 0, 120000)
	h := make(map[string]struct{})
	for scanner.Scan() {
		line := scanner.Text()
		fact := createFact(line)
		if fact != nil {
			if _, ok := h[fact.Hanzi]; ok {
				// TODO use traditional as well?
				// TODO gob these together
				continue
			}
			facts = append(facts, *fact)
			h[fact.Hanzi] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return facts
}

func createFact(line string) *database.Fact {
	if !hasDictFormat(line) {
		return nil
	}
	split1 := strings.SplitN(line, "[", 2)
	mandarin := strings.Split(strings.TrimSpace(split1[0]), " ")
	_, simplified := mandarin[0], mandarin[1]
	split2 := strings.SplitN(split1[1], "]", 2)
	pinyin, english := split2[0], split2[1]
	if isProperNoun(pinyin) {
		return nil
	}
	return &database.Fact{simplified, pinyin, english}
}

func isProperNoun(pinyin string) bool {
	r := []rune(pinyin)
	return unicode.IsUpper(r[0])
}

func hasDictFormat(entry string) bool {
	if strings.HasPrefix(entry, "#") {
		// is a comment
		return false
	}
	if !strings.ContainsRune(entry, '[') {
		return false
	}
	if !strings.ContainsRune(entry, ']') {
		return false
	}
	return true
}

func ParseVocabList(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	words := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}
	return words
}
