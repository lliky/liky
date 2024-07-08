package union

import (
	"fmt"
	"testing"
)

func TestNewSet(t *testing.T) {
	arr := []string{"A", "B", "C", "D", "E"}
	list := make([]interface{}, len(arr))
	for i, v := range arr {
		list[i] = v
	}
	set := NewSet(list)
	fmt.Println(set.IsSameSet("A", "B"))
	set.Union("A", "B")
	fmt.Println(set.IsSameSet("A", "B"))
	set.Union("A", "C")
	fmt.Println(set.IsSameSet("B", "C"))
	fmt.Println(set.IsSameSet("C", "D"))
}
