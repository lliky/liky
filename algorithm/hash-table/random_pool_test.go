package hash_table

import (
	"fmt"
	"testing"
)

func TestNewPool(t *testing.T) {
	p := NewPool()
	p.Insert(1)
	p.Insert(2)
	p.Insert(3)
	p.Delete(1)

	var r1, r2, r3 int
	for i := 0; i < 1<<20; i++ {
		switch p.GetRandom() {
		case 1:
			r1++
		case 2:
			r2++
		case 3:
			r3++
		default:
			fmt.Println("error")
		}
	}
	fmt.Println("r1: ", r1)
	fmt.Println("r2: ", r2)
	fmt.Println("r3: ", r3)
}
