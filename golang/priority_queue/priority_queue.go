package priority_queue

type PriorityQueue []int

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i] < p[i]
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PriorityQueue) Push(x any) {
	p = append(p, x.(int))
}

func (p PriorityQueue) Pop() any {
	old := p
	n := len(old)
	x := old[n-1]
	p = old[0 : n-1]
	return x
}
