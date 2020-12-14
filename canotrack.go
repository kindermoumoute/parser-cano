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
	Descriptions []Description
}

type Description struct {
	Content    []string
	Place      string
	Date       string
	Event      string
	IsSequence bool
	Timer      time.Duration
}

func (t *CanoTrack) ToSRT(sampleLevel int) *SRT {
	srt := &SRT{}
	sequences := []*Caption(nil)

	for _, desc := range t.Descriptions {
		// Add content
		timer, content := desc.Timer, []string(nil)
		switch {
		case desc.IsSequence && sampleLevel >= 3:
			content = desc.Content
		case desc.IsSequence && sampleLevel == 2:
			if len(desc.Content) > 1 {
				content = desc.Content[:1]
			}
		case !desc.IsSequence && sampleLevel == 4:
			content = desc.Content
		case !desc.IsSequence && sampleLevel == 3:
			if len(desc.Content) > 1 {
				content = desc.Content[:1]
			}
		}
		if len(content) != 0 {
			srt.Captions = append(srt.Captions, newCaption(timer+time.Millisecond, content))
		}

		// Add Metadata
		metadataContent := []string(nil)
		if desc.Place != "" {
			metadataContent = append(metadataContent, fmt.Sprintf(`{\an9}<u>Lieu</u>: %s`, desc.Place))
		}
		if desc.Event != "" {
			metadataContent = append(metadataContent, fmt.Sprintf(`{\an8}<u>Evênement</u>: %s`, desc.Event))
		}
		if desc.Date != "" {
			metadataContent = append(metadataContent, fmt.Sprintf(`{\an7}<u>Date</u>: %s`, desc.Date))
		}
		if len(metadataContent) != 0 {
			sequences = append(sequences, newCaption(timer, metadataContent))
		}
	}

	srt.sort()
	srt.removeOverLap()

	srt.Captions = append(srt.Captions, sequences...)
	srt.sort()

	return srt
}

func newCaption(timer time.Duration, content []string) *Caption {
	return &Caption{
		Start: time.Time{}.Add(timer),
		End:   time.Time{}.Add(timer + defaultCaptionDuration),
		Text:  content,
	}
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
			caption.End = s.Captions[i+1].Start.Add(-time.Millisecond)
		}
	}
}
