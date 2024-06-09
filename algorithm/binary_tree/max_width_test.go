package binary_tree

import (
	"fmt"
	"testing"
)

func TestMaxWidth(t *testing.T) {
	head := NewBinaryTree()
	fmt.Println(MaxWidth(head))
}
