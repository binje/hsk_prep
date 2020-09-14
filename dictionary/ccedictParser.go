package dictionary

import (
	"bufio"
	"log"
	"os"
	"strings"

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
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.ContainsRune(line, '[') {
			continue
		}
		if !strings.ContainsRune(line, ']') {
			continue
		}
		split1 := strings.SplitN(line, "[", 2)
		mandarin := strings.Split(strings.TrimSpace(split1[0]), " ")
		_, simplified := mandarin[0], mandarin[1]
		split2 := strings.SplitN(split1[1], "]", 2)
		pinyin, english := split2[0], split2[1]
		facts = append(facts, database.Fact{simplified, pinyin, english})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return facts
}
