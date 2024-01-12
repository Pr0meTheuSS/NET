package webnodes

import (
	"net"
	"sync"
)

type queuedMessage struct {
	to   *net.UDPAddr
	data []byte
}

var queue = messageQueue{
	messages: map[int64]queuedMessage{},
	mtx:      &sync.Mutex{},
}

func (q *messageQueue) delete(sequence int64) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	delete(q.messages, sequence)
}

func (q *messageQueue) add(sequence int64, data []byte, to *net.UDPAddr) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	queue.messages[sequence] = queuedMessage{
		to:   to,
		data: data,
	}
}

type messageQueue struct {
	messages map[int64]queuedMessage
	mtx      *sync.Mutex
}
