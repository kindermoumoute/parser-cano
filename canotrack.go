package main

import (
	"fmt"
	"sort"
	"time"
)

type CanoTrack struct {
	FileName     string
	Length       time.Duration
	Synopsis     string
	Sequences    []Sequence
	Descriptions []Description
}

type Sequence struct {
	Date        string
	Place       string
	Event       string
	People      string
	Description Description
}

type Description struct {
	Content []string
	Timer   time.Duration
}

func (t *CanoTrack) ToSRT(sampleLevel int) *SRT {
	srt := &SRT{}

	if sampleLevel >= 2 {
		for _, seq := range t.Sequences {
			text := seq.Description.Content
			if len(text) > 0 && sampleLevel == 2 {
				text = text[:1]
			}

			if sampleLevel >= 2 && seq.People != "" {
				text = append([]string{seq.People}, text...)
			}
			if text != nil {
				srt.Captions = append(srt.Captions, &Caption{
					Start: time.Time{}.Add(seq.Description.Timer),
					End:   time.Time{}.Add(seq.Description.Timer + defaultCaptionDuration),
					Text:  text,
				})
			}
		}
	}

	if sampleLevel >= 3 {
		for _, desc := range t.Descriptions {
			text := desc.Content
			if sampleLevel == 3 {
				text = text[:1]
			}
			srt.Captions = append(srt.Captions, &Caption{
				Start: time.Time{}.Add(desc.Timer),
				End:   time.Time{}.Add(desc.Timer + defaultCaptionDuration),
				Text:  text,
			})
		}
	}

	srt.sort()
	srt.removeOverLap()

	for _, seq := range t.Sequences {
		text := []string(nil)
		if seq.Place != "" {
			text = append(text, fmt.Sprintf(`{\an7}<font color=#ff0000><u>Lieu</u>: %s<font>`, seq.Place))
		}
		if seq.Event != "" {
			text = append(text, fmt.Sprintf(`{\an4}<font color=#ff0000><u>Evênement</u>: %s<font>`, seq.Event))
		}
		if seq.Date != "" {
			text = append(text, fmt.Sprintf(`{\an9}<font color=#ff0000><u>Date</u>: %s<font>`, seq.Date))
		}

		srt.Captions = append(srt.Captions, &Caption{
			Start: time.Time{}.Add(seq.Description.Timer),
			End:   time.Time{}.Add(seq.Description.Timer + defaultCaptionDuration),
			Text:  text,
		})
	}

	srt.sort()

	return srt
}

func (s *SRT) sort() {
	sort.Slice(s.Captions, func(i, j int) bool {
		return s.Captions[i].Start.Before(s.Captions[j].Start)
	})

	for i, caption := range s.Captions {
		caption.Seq = i
	}
}

func (s *SRT) removeOverLap() {
	for i, caption := range s.Captions {
		if i+1 == len(s.Captions) {
			continue
		}
		if s.Captions[i+1].Start.Before(caption.End) {
			caption.End = s.Captions[i+1].Start
		}
	}
}
