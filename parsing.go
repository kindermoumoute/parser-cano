package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseCanoTrack(rawInput []byte) *CanoTrack {
	inputString := sanitize(string(rawInput))

	track := &CanoTrack{}

	timers, timerIncides, lines := extractTimersFromLines(inputString)

	for i := range timerIncides {
		sectionLines := []string(nil)
		if i+1 < len(timerIncides) {
			sectionLines = lines[timerIncides[i]:timerIncides[i+1]]
		} else {
			sectionLines = lines[timerIncides[i]:]
		}

		switch {
		case i == 0:
			track.FileName, track.Synopsis = parseTrackInfo(lines[:timerIncides[1]])
			track.Length = timers[0]
		case isSequence(sectionLines[0]):
			track.Sequences = append(track.Sequences, parseSequence(sectionLines[1:], timers[i]))
		default:
			track.Descriptions = append(track.Descriptions, parseDescription(sectionLines, timers[i]))
		}
	}
	return track
}

var timerRegexp = regexp.MustCompile(`.*([0-9]{2}:[0-9]{2}).*`)

// extractTimersFromLines find all timers in the raw input and return the following slices:
// timers: all timers found
// timerIndices: line indices of the above timers
// lines: lines from raw input with timers removed
func extractTimersFromLines(rawInput string) ([]time.Duration, []int, []string) {
	lineIndices := make(map[string]int)
	lines := strings.Split(rawInput, "\n")
	for i, line := range lines {
		lineIndices[line] = i
	}
	matches := timerRegexp.FindAllStringSubmatch(rawInput, -1)

	timers := []time.Duration(nil)
	timerIndices := []int(nil)
	for _, lineMatch := range matches {
		lineIndex, exist := lineIndices[lineMatch[0]]
		if !exist {
			fmt.Println("INFO: skipping match on line ", lineMatch)
			continue
		}

		// Convert timer to time.Duration.
		timerParts := strings.Split(lineMatch[1], ":")
		minutes, err := strconv.Atoi(timerParts[0])
		assertNoError(err)
		seconds, err := strconv.Atoi(timerParts[1])
		assertNoError(err)
		timers = append(timers, time.Duration(minutes)*time.Minute+time.Duration(seconds)*time.Second)

		timerIndices = append(timerIndices, lineIndex)
		lines[lineIndex] = sanitizeLine(strings.TrimPrefix(lines[lineIndex], lineMatch[1]))
	}

	return timers, timerIndices, lines
}

func parseTrackInfo(trackInfoLines []string) (string, string) {
	fileName := ""
	synopsisLines := []string(nil)
	for _, line := range trackInfoLines {
		switch {
		case strings.HasPrefix(line, "Fichier:"):
			fileName = strings.TrimPrefix(line, "Fichier:")
		case strings.HasPrefix(line, "Durée:"):
		case strings.HasPrefix(line, "Note boîte:"):
			synopsisLines = append(synopsisLines, strings.TrimPrefix(line, "Note boîte:"))
		default:
			synopsisLines = append(synopsisLines, line)
		}
	}

	return fileName, sanitize(strings.Join(synopsisLines, "\n"))
}

// sanitize does the following sanitizing functions:
// - remove all trailing spaces from lines
// - remove duplicated spaces from str
// - remove all trailing spaces from str
// - remove all empty lines
func sanitize(str string) string {
	lines := []string(nil)
	for _, line := range strings.Split(str, "\n") {
		sanitizedLine := sanitizeLine(line)
		if len(sanitizedLine) != 0 {
			lines = append(lines, sanitizedLine)
		}
	}
	str = strings.Join(lines, "\n")
	str = strings.TrimSpace(str)
	return str
}

func sanitizeLine(line string) string {
	line = strings.TrimSpace(line)
	return strings.Join(strings.Fields(line), " ")
}

// isSequence returns true when first line contains "--- Nouvelle"
func isSequence(sectionFirstLine string) bool {
	return strings.Contains(sectionFirstLine, "--- Nouvelle")
}

func parseSequence(sequenceLines []string, timer time.Duration) Sequence {
	seq := Sequence{
		Description: Description{
			Timer: timer,
		},
	}
	for _, line := range sequenceLines {
		switch {
		case strings.HasPrefix(line, "Date:"):
			seq.Date = sanitizeLine(strings.TrimPrefix(line, "Date:"))
		case strings.HasPrefix(line, "Lieu:"):
			seq.Place = sanitizeLine(strings.TrimPrefix(line, "Lieu:"))
		case strings.HasPrefix(line, "Evênement:"):
			seq.Event = sanitizeLine(strings.TrimPrefix(line, "Evênement:"))
		case strings.HasPrefix(line, "Evenement:"):
			seq.Event = sanitizeLine(strings.TrimPrefix(line, "Evenement:"))
		case strings.HasPrefix(line, "Présents:"):
			seq.People = sanitizeLine(strings.TrimPrefix(line, "Présents:"))
		default:
			seq.Description.Content = append(seq.Description.Content, line)
		}
	}

	return seq
}

func parseDescription(descriptionLines []string, timer time.Duration) Description {
	return Description{
		Content: descriptionLines,
		Timer:   timer,
	}
}
