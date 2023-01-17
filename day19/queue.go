package main

// Number a type placeholder is used to specify either an integer or float64
type Number interface {
	int | float64
}

// Item can be any type of value and a number
type Item[T any, N Number] struct {
	item  T
	value N
}

// Queue holds Items of type T
// A value function takes type T and returns a value N for that item
type Queue[T any, N Number] struct {
	items []Item[T, N]
	value func(item T) N
}

// Create a NewPriorityQueue object of type T
// Takes a function val as a parameter which is used to get the priority values for each item in the queue
// The priority queue is sorted from high priority items to low
func NewPriorityQueue[T any, N Number](val func(item T) N) *Queue[T, N] {
	h := new(Queue[T, N])
	h.value = val
	return h
}

// Create a NewQueue object with type T
func NewQueue[T any, N Number]() *Queue[T, N] {
	h := new(Queue[T, N])
	return h
}

// Enqueue is used to add item of type T into the queue h
// For a simple queue the item is appended to the slice of items
// For a priority queue a value will be used to calculate that items priority in the queue
// Items with higher priority are added in their priority order
func (h *Queue[T, N]) Enqueue(item T) {
	if h.value == nil {
		h.items = append(h.items, Item[T, N]{item: item})
	} else {
		v := h.value(item)
		idx := -1
		for i := 0; i < len(h.items); i++ {
			r := h.items[i].value

			if v > r {
				idx = i
				break
			}
		}

		i := Item[T, N]{item: item, value: v}
		if idx == -1 {
			h.items = append(h.items, i)
		} else {
			h.items = append(h.items[:idx], append([]Item[T, N]{i}, h.items[idx:]...)...)
		}
	}
}

// Dequeue is used to remove the first item from the Queue
// Either the item that has been there the longest for a normal Queue or the highest priority item for a priority Queue
// If no items are in the Queue when Dequeue is called then it will return the zero value of type T
func (h *Queue[T, N]) Dequeue() T {
	var item Item[T, N]
	if len(h.items) > 0 {
		item = h.items[0]
		h.items = h.items[1:]
	}
	return item.item
}

// Any returns true is there are Any items in Queue h otherwise false if it is empty
func (h *Queue[T, N]) Any() bool {
	return h.Len() > 0
}

// Len returns the number of items in Queue h
func (h *Queue[T, N]) Len() int {
	return len(h.items)
}
