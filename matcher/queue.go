package matcher

// queue holds <capacity> number of elements
type queue struct {
	elms     []byte
	capacity int
}

// newQueue creates and returns a new queue with given capacity
func newQueue(capacity int) *queue {
	return &queue{
		elms:     make([]byte, 0, capacity),
		capacity: capacity,
	}
}

// add adds a new element to queue. If new element exceeds the capacity,
// it removes the first element from the front and returns it
func (q *queue) add(elm byte) (b byte, overflow bool) {
	if q.capacity == 0 {
		return elm, true
	}
	if len(q.elms) == q.capacity {
		b = q.pop()
		overflow = true
	}
	q.elms = append(q.elms, elm)
	return b, overflow
}

func (q *queue) getElements() []byte {
	return q.elms
}

// pop removes the first element from the front and returns it
func (q *queue) pop() (b byte) {
	b = q.elms[0]
	q.elms = q.elms[1:]
	return b
}
