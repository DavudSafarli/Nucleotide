package matcher

import (
	"bufio"
	"io"
)

// SequenceMatcherOptions are the options for creating a new SequenceMatcher
type SequenceMatcherOptions struct {
	precedingCharCount  int
	succeedingCharCount int
	sequence            []byte
	eos                 rune
}

// MatchResult represents the result of one found match
type MatchResult struct {
	preceding  string
	match      string
	succeeding string
}

// SequenceMatcher is used for finding sequences in a stream with left and right contexts.
type SequenceMatcher struct {
	reader  io.Reader
	options SequenceMatcherOptions

	ch chan MatchResult

	// SequenceMatcher holds a queue for each of Succeeding, Current and Preceding elements.
	preQueue      *queue
	matchingQueue *queue
	sucQueue      *queue
}

// NewSequenceMatcher returns a new SequenceMatcher for finding sequences in a stream
func NewSequenceMatcher(reader io.Reader, options SequenceMatcherOptions) SequenceMatcher {
	return SequenceMatcher{
		reader:        reader,
		options:       options,
		ch:            make(chan MatchResult, 1),
		preQueue:      newQueue(options.precedingCharCount),
		matchingQueue: newQueue(len(options.sequence)),
		sucQueue:      newQueue(options.succeedingCharCount),
	}
}

// sendMatch sends current Match to the channel
func (m SequenceMatcher) sendMatch() {
	m.ch <- MatchResult{
		preceding:  string(m.preQueue.getElements()),
		match:      string(m.matchingQueue.getElements()),
		succeeding: string(m.sucQueue.getElements()),
	}
}

// addByte adds byte element to SucceedingQueue.
// Overflowed element from the SucceedingQueue will be added to MatchingQueue.
// Overflowed element of MatchingQueue will also be added to PrecedingQueue.
// Overflowed element of PrecedingQueue will be thrown away
func (m SequenceMatcher) addByte(b byte) {
	b, overflowed := m.sucQueue.add(b)
	if !overflowed {
		return
	}
	b, overflowed = m.matchingQueue.add(b)
	if !overflowed {
		return
	}
	_, _ = m.preQueue.add(b)
}

// pop pops one element from SucceedingQueue adds it to MatchingQueue.
// Overflowed element of MatchingQueue will be added to PrecedingQueue.
// Overflowed element of PrecedingQueue will be thrown away
func (m SequenceMatcher) pop() {
	b := m.sucQueue.pop()
	b, _ = m.matchingQueue.add(b)
	m.preQueue.add(b)
}

// isMatch checks if current elements in the MatchingQueue is a valid match for wanted sequence
func (m SequenceMatcher) isMatch() bool {
	sequence := m.matchingQueue.getElements()
	if len(sequence) != len(m.options.sequence) {
		return false
	}

	for i := 0; i < len(sequence); i++ {
		if sequence[i] != m.options.sequence[i] {
			return false
		}
	}
	return true
}

// Run starts the matching process
func (m SequenceMatcher) Run() <-chan MatchResult {
	go m.readStreamAndMatchSequences()
	return m.ch
}

// readStreamAndMatchSequences reads from stream rune-by-rune until encountering EOS or EOF.
func (m SequenceMatcher) readStreamAndMatchSequences() {
	defer close(m.ch)
	reader := bufio.NewReader(m.reader)

	// Add new elements to queue one-by-one
	// Check if current sequence is a Match
	for {
		char, size, err := reader.ReadRune()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF || char == m.options.eos {
			break
		}
		// if read rune is not 1-byte long, it means there are non-ascii characters in the stream other than EOS
		if size != 1 {
			panic("encountered non ascii character which is not EOS either. Program only supports ascii charset for stream except EOS.")
		}

		m.addByte(byte(char))
		if m.isMatch() {
			m.sendMatch()
		}
	}
	// turn the wheel for number of elements in the SucceedingQueue too
	// and check possible matches
	l := len(m.sucQueue.getElements())
	for i := 0; i < l; i++ {
		m.pop()
		if m.isMatch() {
			m.sendMatch()
		}
	}
}
