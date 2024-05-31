package algorithm

func Fib(n int) int {
	if n <= 1 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func Fib2(n int, prev, next int) int {
	if n == 0 {
		return prev
	}
	n--
	return Fib2(n, next, next+prev)
}
