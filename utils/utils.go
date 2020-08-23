package utils

import (
	"strings"
)

var accents = map[rune]map[int]rune{
	'a': map[int]rune{
		'1': '\u0101',
		'2': '\u00e1',
		'3': '\u01ce',
		'4': '\u00e0'},
	'e': map[int]rune{
		'1': '\u0113',
		'2': '\u00e9',
		'3': '\u011b',
		'4': '\u00e8'},
	'i': map[int]rune{
		'1': '\u012b',
		'2': '\u00ed',
		'3': '\u01d0',
		'4': '\u00ec'},
	'o': map[int]rune{
		'1': '\u014d',
		'2': '\u00f3',
		'3': '\u01d2',
		'4': '\u00f2'},
	'u': map[int]rune{
		'1': '\u016b',
		'2': '\u00fa',
		'3': '\u01d4',
		'4': '\u00f9'},
	'ü': map[int]rune{
		'1': 'ǖ',
		'2': 'ǘ',
		'3': 'ǚ',
		'4': 'ǜ'},
}

type pInfo struct {
	c    string
	tone string
}

var pinyin = map[rune]pInfo{
	'\u0101': pInfo{"a", "1"},
	'\u01e1': pInfo{"a", "2"},
	'\u01ce': pInfo{"a", "3"},
	'\u01c0': pInfo{"a", "4"},
	'\u0113': pInfo{"e", "1"},
	'\u00e9': pInfo{"e", "2"},
	'\u011b': pInfo{"e", "3"},
	'\u00e8': pInfo{"e", "4"},
	'\u012b': pInfo{"i", "1"},
	'\u00ed': pInfo{"i", "2"},
	'\u01d0': pInfo{"i", "3"},
	'\u00ec': pInfo{"i", "4"},
	'\u014d': pInfo{"o", "1"},
	'\u00f3': pInfo{"o", "2"},
	'\u01d2': pInfo{"o", "3"},
	'\u00f2': pInfo{"o", "4"},
	'\u016b': pInfo{"u", "1"},
	'\u00fa': pInfo{"u", "2"},
	'\u01d4': pInfo{"u", "3"},
	'\u00f9': pInfo{"u", "4"},
	'ǖ':      pInfo{"ü", "1"},
	'ǘ':      pInfo{"ü", "2"},
	'ǚ':      pInfo{"ü", "3"},
	'ǜ':      pInfo{"ü", "4"},
}

func Typable2PinyinByLine(line string) string {
	words := strings.Fields(line)
	for n, word := range words {
		words[n] = Typable2Pinyin(word)
	}
	return strings.Join(words, ` `)
}

func Typable2Pinyin(word string) string {
	word = strings.ReplaceAll(word, "u:", "ü")
	chars := []rune(word)
	tone := chars[len(chars)-1]
	if tone == '5' {
		return string(chars[0 : len(chars)-1])
	} else if tone >= '1' || tone <= '4' {
		word = string(chars[0 : len(chars)-1])
		word = addAccent(word, int(tone))
	}
	return word
}

func addAccent(word string, tone int) string {
	n := strings.Index(word, "a")
	if n != -1 {
		return strings.Replace(word, "a", string(accents['a'][tone]), 1)
	}
	n = strings.Index(word, "e")
	if n != -1 {
		return strings.Replace(word, "e", string(accents['e'][tone]), 1)
	}
	n = strings.Index(word, "ou")
	if n != -1 {
		return strings.Replace(word, "o", string(accents['o'][tone]), 1)
	}
	chars := []rune(word)
	for i := len(chars) - 1; i >= 0; i-- {
		switch chars[i] {
		case 'i':
			chars[i] = accents['i'][tone]
			return string(chars)
		case 'o':
			chars[i] = accents['o'][tone]
			return string(chars)
		case 'u':
			chars[i] = accents['u'][tone]
			return string(chars)
		case 'ü':
			chars[i] = accents['ü'][tone]
			return string(chars)
		}
	}
	return word
}

func TypablePinyinByLine(line string) string {
	words := strings.Fields(line)
	for n, word := range words {
		words[n] = TypablePinyin(word)
	}
	return strings.Join(words, ` `)
}

func TypablePinyin(word string) string {
	for _, c := range word {
		if a, ok := pinyin[c]; ok {
			word = strings.Replace(word, string(c), a.c, 1)
			return word + a.tone
		}
	}
	return word + "5"
}

/*
var prefixes = []string{
	"b",
	"p",
	"m",
	"f",
	"d",
	"t",
	"n",
	"l",
	"g",
	"k",
	"h",
	"j",
	"q",
	"x",
	"zh",
	"ch",
	"sh",
	"r",
	"z",
	"c",
	"s"}

var suffixes = []string{
	"a",
	"o",
	"ie",
	"ai",
	"ei",
	"ao",
	"ou",
	"an",
	"en",
	"ang",
	"eng",
	"er",
	"i",
	"u"}

var words = []string{}

func typablePinyin(input string) string {
	if len(words) == 0 {
		for _, i := range prefixes {
			for _, j := range suffixes {
				words = append(words, fmt.Sprintf("%s%s", i, j))
				fmt.Printf("%v%v,\n", i, j)
			}
		}
	}
	return ""
}
*/
