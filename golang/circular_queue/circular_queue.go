package circular_queue

type CircularQueue struct {
	data  []interface{}
	size  int
	front int
	rear  int
}

func NewCircularQueue(size int) *CircularQueue {
	return &CircularQueue{
		data:  make([]interface{}, size),
		size:  size,
		front: 0,
		rear:  0,
	}
}

func (c *CircularQueue) Enqueue(value interface{}) bool {
	if c.rear == c.front+c.size {
		return false
	}
	c.data[c.rear%c.size] = value
	c.rear++
	return true
}

func (c *CircularQueue) Dequeue() (interface{}, bool) {
	if c.rear == c.front {
		return nil, false
	}
	value := c.data[c.front%c.size]
	c.front++
	return value, true
}

func (c *CircularQueue) Len() int {
	return c.rear - c.front
}
