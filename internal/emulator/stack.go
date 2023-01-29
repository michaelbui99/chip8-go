package emulator

import "errors"

type Stack struct {
	stack    []uint16
	capacity int
}

func NewStack(capacity int) *Stack {
	return &Stack{
		stack:    []uint16{},
		capacity: capacity,
	}
}

func (stack *Stack) Push(element uint16) error {
	if stack.Size() >= stack.capacity {
		return errors.New("stack has already reached max capacity")
	}

	stack.stack = append(stack.stack, element)
	return nil
}

func (stack *Stack) Pop() (uint16, error) {
	if stack.Empty() {
		return 0, errors.New("empty stack")
	}

	element := stack.stack[stack.Size()-1]
	stack.stack = stack.stack[:stack.Size()-1]

	return element, nil
}

func (stack *Stack) Peek() (uint16, error) {
	if stack.Empty() {
		return 0, errors.New("empty stack")
	}

	return stack.stack[stack.Size()-1], nil
}

func (stack *Stack) Size() int {
	return len(stack.stack)
}

func (stack *Stack) Empty() bool {
	return stack.Size() > 0
}
