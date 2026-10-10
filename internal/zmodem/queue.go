package zmodem

import (
	"errors"
	"io"
	"sync"
)

var errQueueFull = errors.New("zmodem: proto buffer full")

// byteQueue is a bounded buffer. tryWrite never blocks.
type byteQueue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	buf    []byte
	max    int
	closed bool
}

func newByteQueue(max int) *byteQueue {
	q := &byteQueue{max: max}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *byteQueue) tryWrite(p []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return io.ErrClosedPipe
	}
	if len(q.buf)+len(p) > q.max {
		return errQueueFull
	}
	q.buf = append(q.buf, p...)
	q.cond.Signal()
	return nil
}

func (q *byteQueue) Read(p []byte) (int, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.buf) == 0 && !q.closed {
		q.cond.Wait()
	}
	if len(q.buf) == 0 {
		return 0, io.EOF
	}
	n := copy(p, q.buf)
	q.buf = q.buf[n:]
	if len(q.buf) == 0 {
		q.buf = nil
	}
	return n, nil
}

func (q *byteQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	q.cond.Broadcast()
}
