package main

import (
	"bytes"
	"testing"
)

//type SequenceMatcher struct {
//	precedingCharCount int
//	succeedingCharCount int
//	sequence string
//	r io.Reader
//}
//
//func NewSequenceMatcher(x, y int, sequence string) SequenceMatcher {
//	return SequenceMatcher{
//		precedingCharCount:  x,
//		succeedingCharCount: y,
//		sequence:            sequence,
//	}
//}

func TestSequenceMatcher(t *testing.T) {
	type inputConfig struct {
		x     int
		y     int
		T     string
		input string
	}
	tests := []struct {
		config inputConfig
		ans    MatchResult
	}{
		{
			config: inputConfig{
				x:     2,
				y:     2,
				T:     "AAA",
				input: "SAAALAM",
			},
			ans: MatchResult{
				match:     "AAA",
			},
		},
	}
	for _, tt := range tests {
		buff := &bytes.Buffer{}
		matcher := NewSequenceMatcher(buff, tt.config.x, tt.config.y, tt.config.T)
		buff.WriteString(tt.config.input)
		var matches <-chan MatchResult = matcher.Run()
		if tt.ans.match != (<-matches).match {
			t.Error("Failed")
		}
	}
}
