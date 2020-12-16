package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	parser_cano "github.com/kindermoumoute/parser-cano"
)

func main() {
	// go run . 002C.txt 003C.txt 004C.txt 005C.txt 007C.txt 009C.txt

	// go run . --stdout --sample-level=4 009C.txt

	sampleLevel := flag.Int("sample-level", 1, "verbosity of the sample (1 to 4)")
	displayOnStdout := flag.Bool("stdout", false, "display output on stdout")

	flag.Parse()

	for _, file := range flag.Args() {
		rawInput, err := ioutil.ReadFile(file)
		assertNoError(err)

		if !*displayOnStdout {
			fmt.Println("parsing file", file)
		}

		track := parser_cano.ParseCanoTrack(rawInput)
		srtContent := track.ToSRT(*sampleLevel).ToSRT()
		if *displayOnStdout {
			fmt.Println(srtContent)
		} else {
			outputFileName := fmt.Sprintf("%s-%d.srt", track.FileName, *sampleLevel)
			fmt.Println("writing file", outputFileName)
			assertNoError(ioutil.WriteFile(outputFileName, []byte(srtContent), 0644))
		}
	}
}

func assertNoError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
