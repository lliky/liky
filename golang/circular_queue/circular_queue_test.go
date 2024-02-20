package circular_queue

import "testing"

func TestNewCircularQueue(t *testing.T) {
	results := []struct {
		op       string
		data     []int
		size     int
		deCount  int
		expected int
	}{
		{
			op:       "fullQueue",
			data:     []int{1, 2, 3, 4, 5, 6},
			size:     5,
			deCount:  0,
			expected: 5,
		},
		{
			op:       "emptyQueue",
			data:     []int{1, 2, 3, 4},
			size:     5,
			deCount:  4,
			expected: 0,
		},
		{
			op:       "enqueue",
			data:     []int{1, 2, 3, 4},
			size:     5,
			deCount:  0,
			expected: 4,
		},
	}

	for _, result := range results {
		t.Run(result.op, func(t *testing.T) {
			c := NewCircularQueue(result.size)
			for _, value := range result.data {
				c.Enqueue(value)
			}
			for i := 0; i < result.deCount; i++ {
				value, _ := c.Dequeue()
				if value != result.data[i] {
					t.Errorf("Expected %d, got %d", result.data[i], value)
				}
			}
			if c.Len() != result.expected {
				t.Errorf("The length of Queue expected %d, got %d", result.expected, c.Len())
			}
		})
	}
}
