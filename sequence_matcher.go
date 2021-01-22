package main

import (
	"io"
	"reflect"
)

type SequenceMatcher struct {
	precedingCharCount  int
	succeedingCharCount int
	sequence            []byte
	reader              io.Reader
}

func NewSequenceMatcher(reader io.Reader, x, y int, sequence string) SequenceMatcher {
	return SequenceMatcher{
		precedingCharCount:  x,
		succeedingCharCount: y,
		sequence:            []byte(sequence),
		reader:              reader,
	}
}

func (sm SequenceMatcher) Run() <-chan MatchResult {
	a := make(chan MatchResult, 1)
	go sm.listenStreamAndMatchSequence(a)
	return a
}

func (sm SequenceMatcher) listenStreamAndMatchSequence(ch chan<- MatchResult) {
	window := make([]byte, len(sm.sequence), len(sm.sequence))
	for {
		newChar := make([]byte, 1)
		_, err := sm.reader.Read(newChar)
		if err == io.EOF {
			close(ch)
			break
		}
		if err != nil {
			close(ch)
			panic(err)
		}

		window = window[1:]
		window = append(window, newChar[0])
		if !reflect.DeepEqual(window, sm.sequence) {
			continue
		}

		ch <- MatchResult{
			match: string(sm.sequence),
		}
	}
}

type MatchResult struct {
	match     string
	count     int
	preceding string
	succeding string
}
