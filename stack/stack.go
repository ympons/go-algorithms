package stack

// Stack interface
type Stack interface {
	Push(int)
	Pop() (int, bool)
	IsEmpty() bool
}

type snode struct {
	value int
	next  *snode
}

type stack struct {
	top *snode
	n   int
}

// NewStack creates a new Stack
func NewStack() Stack {
	return &stack{}
}

func (s *stack) Push(v int) {
	s.top = &snode{v, s.top}
	s.n++
}

func (s *stack) Pop() (value int, ok bool) {
	if ok = s.n > 0; ok {
		value, s.top = s.top.value, s.top.next
		s.n--
	}
	return
}

func (s stack) IsEmpty() bool {
	return s.n == 0
}
