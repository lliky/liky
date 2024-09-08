package bitwise_operation

import (
	"fmt"
	"math"
	"testing"
)

func TestGetMax1(t *testing.T) {
	a := math.MinInt
	b := math.MaxInt
	fmt.Println(GetMax1(a, b)) // error, because of overflow
}

func TestGetMax2(t *testing.T) {
	fmt.Println(GetMax2(1, 5))
	fmt.Println(GetMax2(math.MinInt, math.MaxInt))
}
