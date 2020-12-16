package parser_cano

import (
	"sort"
	"strings"
	"time"
)

const (
	sequenceCaptionDuration    = 10 * time.Second
	minimumInfoCaptionDuration = 3 * time.Second
	charPerSecond              = 15
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
			if len(desc.Content) >= 1 {
				content = desc.Content[:1]
			}
		case !desc.IsSequence && sampleLevel == 4:
			content = desc.Content
		case !desc.IsSequence && sampleLevel == 3:
			if len(desc.Content) >= 1 {
				content = desc.Content[:1]
			}
		}
		if len(content) != 0 {
			charCount := len(strings.Join(content, ""))
			infoDuration := time.Duration(charCount/charPerSecond)*time.Second + minimumInfoCaptionDuration
			srt.Captions = append(srt.Captions, newCaption(timer+time.Millisecond, content, infoDuration))
		}

		// Add Metadata
		if desc.Place != "" {
			sequences = append(sequences, newCaption(timer, []string{`{\an9}` + desc.Place}, sequenceCaptionDuration))
		}
		if desc.Event != "" {
			sequences = append(sequences, newCaption(timer, []string{`{\an8}` + desc.Event}, sequenceCaptionDuration))
		}
		if desc.Date != "" {
			sequences = append(sequences, newCaption(timer, []string{`{\an7}` + desc.Date}, sequenceCaptionDuration))
		}
	}

	srt.sort()
	srt.removeOverLap()

	srt.Captions = append(srt.Captions, sequences...)
	srt.sort()

	return srt
}

func newCaption(timer time.Duration, content []string, duration time.Duration) *Caption {
	return &Caption{
		Start: time.Time{}.Add(timer),
		End:   time.Time{}.Add(timer + duration),
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
