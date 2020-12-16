package parser_cano

import (
	"strings"
	"time"
)

func (t *CanoTrack) ToIndex(keyword string) string {
	results := []string{}
	//fmt.Println(t.Synopsis)
	if strings.Contains(t.Synopsis, keyword) {
		results = append(results, newResult(t.FileName, 0, "", "", "", t.Synopsis))
	}

	lastPlace, lastDate, lastEvent := "", "", ""
	for _, desc := range t.Descriptions {
		if desc.Place != "" {
			lastPlace = desc.Place
		}
		if desc.Date != "" {
			lastDate = desc.Date
		}
		if desc.Event != "" {
			lastEvent = desc.Event
		}
		if searchKeyword(desc.Place, keyword) ||
			searchKeyword(desc.Date, keyword) ||
			searchKeyword(desc.Event, keyword) ||
			searchKeyword(strings.Join(desc.Content, "\n"), keyword) {
			results = append(results, newResult(t.FileName, desc.Timer, lastDate, lastPlace, lastEvent, strings.Join(desc.Content, "\n")))
		}
	}

	return strings.Join(results, "\n")
}

func searchKeyword(s, keyword string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(keyword))

}

func newResult(filename string, timer time.Duration, lastDate string, lastPlace, lastEvent string, fullLine string) string {
	return strings.Join([]string{
		filename + ".mp4",
		time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC).Add(timer).Format("04:05"),
		lastPlace,
		lastDate,
		lastEvent,
		strings.ReplaceAll(strings.ReplaceAll(fullLine, ",", " "), "\n", " - "),
	}, ",")
}
