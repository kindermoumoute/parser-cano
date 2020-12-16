package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	parser_cano "github.com/kindermoumoute/parser-cano"
)

func main() {
	// go run main.go 002C.txt 003C.txt 004C.txt 005C.txt 007C.txt 009C.txt

	flag.Parse()

	for _, file := range flag.Args() {
		rawInput, err := ioutil.ReadFile(file)
		assertNoError(err)

		track := parser_cano.ParseCanoTrack(rawInput)

		words := []string(nil)
		wordsCount := make(map[string]int)
		validContent := track.Synopsis
		for _, desc := range track.Descriptions {
			validContent += "\n" + desc.Event
			validContent += "\n" + desc.Place
			validContent += "\n" + desc.Date
			validContent += "\n" + strings.Join(desc.Content, "\n")
		}
		for _, word := range strings.Fields(validContent) {
			word = cleanWord(word)
			if shouldSkipWord(word) {
				continue
			}
			_, exist := wordsCount[word]
			if !exist {
				words = append(words, word)
			}
			wordsCount[word]++
		}

		sort.Slice(words, func(i, j int) bool {
			return wordsCount[words[i]] > wordsCount[words[j]]
		})

		for _, word := range words {
			if wordsCount[word] > 1 {
				fmt.Printf("%s,%d\n", word, wordsCount[word])
			}
		}
	}
}

func shouldSkipWord(word string) bool {
	if len(word) <= 3 {
		return true
	}
	switch word {
	case "puis", "avec":
		return true
	}
	return false
}

func cleanWord(word string) string {
	return strings.Trim(strings.ToLower(word), "()")
}

func assertNoError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
