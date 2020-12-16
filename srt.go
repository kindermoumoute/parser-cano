package parser_cano

import (
	"fmt"
	"strings"
	"time"
)

// Eol is the end of line characters to use when writing .srt data
var eol = "\n"

type SRT struct {
	Captions []*Caption
}

// AsSRT renders the sub in .srt format
func (s *SRT) ToSRT() string {
	res := ""
	for _, sub := range s.Captions {
		res += sub.ToSRT()
	}
	return res
}

// Caption represents one subtitle block
type Caption struct {
	Seq   int
	Start time.Time
	End   time.Time
	Text  []string
}

// AsSRT renders the caption as srt
func (c Caption) ToSRT() string {
	res := fmt.Sprintf("%d", c.Seq) + eol +
		TimeSRT(c.Start) + " --> " + TimeSRT(c.End) + eol
	for _, line := range c.Text {
		res += line + eol
	}
	return res + eol
}

// TimeSRT renders a timestamp for use in .srt
func TimeSRT(t time.Time) string {
	res := t.Format("15:04:05.000")
	return strings.Replace(res, ".", ",", 1)
}
