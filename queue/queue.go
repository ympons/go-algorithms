package queue

// Queue interface
type Queue interface {
	Enqueue(int)
	Dequeue() (int, bool)
	IsEmpty() bool
}

type qnode struct {
	value int
	next  *qnode
}

type queue struct {
	head *qnode
	tail *qnode
	n    int
}

// NewQueue creates a new Queue
func NewQueue() Queue {
	return &queue{}
}

func (q *queue) Enqueue(v int) {
	node := &qnode{value: v}
	if q.tail == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.n++
}

func (q *queue) Dequeue() (int, bool) {
	if q.head == nil {
		return 0, false
	}

	node := q.head
	q.head = node.next
	if q.head == nil {
		q.tail = nil
	}
	q.n--

	return node.value, true
}

func (q queue) IsEmpty() bool {
	return q.n == 0
}
