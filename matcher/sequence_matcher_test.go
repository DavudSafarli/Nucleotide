package matcher

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func BenchSequenceMatcher(input string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		buff := &bytes.Buffer{}
		buff.WriteString(input)
		matcher := NewSequenceMatcher(buff, SequenceMatcherOptions{
			precedingCharCount:  5,
			sequence:            []byte("AGTA"),
			succeedingCharCount: 7,
			eos:                 'ε',
		})
		var matchesChan <-chan MatchResult = matcher.Run()
	upper:
		for {
			select {
			case _, isOpen := <-matchesChan:
				if !isOpen {
					break upper
				}
			}
		}
	}
}

func BenchmarkSequenceMatcherX1(b *testing.B) {
	BenchSequenceMatcher("AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε", b)
}
func BenchmarkSequenceMatcherX2(b *testing.B) {
	BenchSequenceMatcher("AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAAAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε", b)
}
func BenchmarkSequenceMatcherX4(b *testing.B) {
	BenchSequenceMatcher("AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAAAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAAAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAAAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε", b)
}

func TestSequenceMatcher(t *testing.T) {
	type inputConfig struct {
		x     int
		y     int
		T     string
		input string
		eos   rune
	}
	tests := []struct {
		config inputConfig
		ans    []MatchResult
	}{
		{
			config: inputConfig{
				x:     1,
				y:     2,
				T:     "ACGT",
				input: "ACACGTCAε",
				eos:   'ε',
			},
			ans: []MatchResult{
				{"C", "ACGT", "CA"},
			},
		},
		{
			config: inputConfig{
				x:     1,
				y:     5,
				T:     "ACGT",
				input: "ACACGTCAε",
				eos:   'ε',
			},
			ans: []MatchResult{
				{"C", "ACGT", "CA"},
			},
		},
		{
			config: inputConfig{
				x:     1,
				y:     1,
				T:     "AAA",
				input: "SAAALAM",
			},
			ans: []MatchResult{
				{"S", "AAA", "L"},
			},
		},
		{
			config: inputConfig{
				x:     2,
				y:     2,
				T:     "AAA",
				input: "SAAALAAAM",
			},
			ans: []MatchResult{
				{"S", "AAA", "LA"},
				{"AL", "AAA", "M"},
			},
		},
		{
			config: inputConfig{
				x:     5,
				y:     7,
				T:     "AGTA",
				input: "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTA",
			},
			ans: []MatchResult{
				{"A", "AGTA", "CGTGCAG"},
				{"CAGTG", "AGTA", "GTAGACC"},
				{"TGAGT", "AGTA", "GACCTGA"},
				{"ATATA", "AGTA", "GCTA"},
			},
		},
		{
			config: inputConfig{
				x:     0,
				y:     0,
				T:     "AGTA",
				input: "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTA",
			},
			ans: []MatchResult{
				{"", "AGTA", ""},
				{"", "AGTA", ""},
				{"", "AGTA", ""},
				{"", "AGTA", ""},
			},
		},
		{
			config: inputConfig{
				x:     0,
				y:     1,
				T:     "AGTA",
				input: "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTA",
			},
			ans: []MatchResult{
				{"", "AGTA", "C"},
				{"", "AGTA", "G"},
				{"", "AGTA", "G"},
				{"", "AGTA", "G"},
			},
		},
		{
			config: inputConfig{
				x:     1,
				y:     0,
				T:     "AGTA",
				input: "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTA",
			},
			ans: []MatchResult{
				{"A", "AGTA", ""},
				{"G", "AGTA", ""},
				{"T", "AGTA", ""},
				{"A", "AGTA", ""},
			},
		},
	}
	for i, tt := range tests {
		t.Run(`test `+strconv.Itoa(i), func(t *testing.T) {
			// arrange
			buff := &bytes.Buffer{}
			matcher := NewSequenceMatcher(buff, SequenceMatcherOptions{
				precedingCharCount:  tt.config.x,
				sequence:            []byte(tt.config.T),
				succeedingCharCount: tt.config.y,
				eos:                 tt.config.eos,
			})

			// act
			buff.WriteString(tt.config.input)
			var matchesChan <-chan MatchResult = matcher.Run()

			matches := make([]MatchResult, 0, len(tt.ans))
			tick := time.Tick(100 * time.Second)

			// put the result matches into an array for "1 second or until channel is closed"
		upper:
			for {
				select {
				case <-tick:
					t.Error("channel was not closed in 1 second")
					t.FailNow()
				case match, isOpen := <-matchesChan:
					if !isOpen {
						break upper
					}
					matches = append(matches, match)
				}
			}
			if len(matches) != len(tt.ans) {
				t.Error("do not have the same length")
				t.Error("\nexpected:", tt.ans, "\nactual:  ", matches)
				t.FailNow()
			}
			for i := 0; i < len(tt.ans); i++ {
				actual := matches[i]
				expected := tt.ans[i]
				if expected != actual {
					t.Error("\nexpected:", tt.ans, "\nactual:  ", matches)
					t.FailNow()
				}
			}
		})
	}
}
