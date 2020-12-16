package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	parser_cano "github.com/kindermoumoute/parser-cano"
)

func main() {
	// go run main.go --keyword=vacance 002C.txt 003C.txt 004C.txt 005C.txt 007C.txt 009C.txt

	keyword := flag.String("keyword", "", "keyword to search")
	flag.Parse()
	if *keyword == "" {
		fmt.Fprintln(os.Stderr, "--keyword flag is required")
		os.Exit(1)
	}

	fmt.Println("File Name,Timer,Date,Place,Event,Description")
	for _, file := range flag.Args() {
		rawInput, err := ioutil.ReadFile(file)
		assertNoError(err)

		track := parser_cano.ParseCanoTrack(rawInput)

		fmt.Println(track.ToIndex(*keyword))
	}
}

func assertNoError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
