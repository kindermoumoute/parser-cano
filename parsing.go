package parser_cano

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	NewSequenceReg = regexp.MustCompile(`.*--- Nouvelle.*`)
	FileNameReg    = regexp.MustCompile(`Fichier`)
	LengthReg      = regexp.MustCompile(`Durée`)
	SynopsisReg    = regexp.MustCompile(`Note boîte`)
	DateReg        = regexp.MustCompile(`Date`)
	PlaceReg       = regexp.MustCompile(`Lieu`)
	PeopleReg      = regexp.MustCompile(`Pr[é|e]{1}sents?`)
	EventReg       = regexp.MustCompile(`Ev[ê|e]{1}nement`)

	timerRegexp = regexp.MustCompile(`.*([0-9]{2}:[0-9]{2}).*`)
)

func ParseCanoTrack(rawInput []byte) *CanoTrack {
	inputString := sanitize(string(rawInput))

	track := &CanoTrack{}

	timers, timerIndices, lines := extractTimersFromLines(inputString)

	for i := range timerIndices {
		sectionLines := []string(nil)
		if i+1 < len(timerIndices) {
			sectionLines = lines[timerIndices[i]:timerIndices[i+1]]
		} else {
			sectionLines = lines[timerIndices[i]:]
		}

		switch {
		case i == 0:
			track.FileName, track.Synopsis = parseTrackInfo(lines[:timerIndices[1]])
			track.Length = timers[0]
		default:
			track.Descriptions = append(track.Descriptions, parseDescription(sectionLines, timers[i]))
		}
	}
	return track
}

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
		if err != nil {
			panic(err)
		}
		seconds, err := strconv.Atoi(timerParts[1])
		if err != nil {
			panic(err)
		}
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
		tmp := strings.SplitN(line, ":", 2)
		content := ""
		if len(tmp) == 2 {
			content = sanitize(tmp[1])
		}
		switch {
		case FileNameReg.MatchString(line):
			fileName = content
		case LengthReg.MatchString(line):
		case SynopsisReg.MatchString(line):
			if content != "" {
				synopsisLines = append(synopsisLines, content)
			}
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

func parseDescription(descriptionLines []string, timer time.Duration) Description {
	desc := Description{
		Timer: timer,
	}
	for _, line := range descriptionLines {
		tmp := strings.SplitN(line, ":", 2)
		content := ""
		if len(tmp) == 2 {
			content = sanitize(tmp[1])
		}
		switch {
		case NewSequenceReg.MatchString(tmp[0]):
			desc.IsSequence = true
		case DateReg.MatchString(tmp[0]):
			if content != "" {
				desc.Date = content
			}
		case PlaceReg.MatchString(tmp[0]):
			if content != "" {
				desc.Place = content
			}
		case EventReg.MatchString(tmp[0]):
			if content != "" {
				desc.Event = content
			}
		case PeopleReg.MatchString(tmp[0]):
			if content != "" {
				desc.Content = append(desc.Content, line)
			}
		default:
			line = sanitizeLine(line)
			if line != "" {
				desc.Content = append(desc.Content, line)
			}
		}
	}

	return desc
}
