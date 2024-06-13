package binary_tree

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBCT(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		f      func() *Node
		expect bool
	}{
		{
			name:   "newValidCompleteTree",
			f:      newValidCompleteTree,
			expect: true,
		},
		{
			name:   "newInvalidCompleteTree",
			f:      newInvalidCompleteTree,
			expect: false,
		},
	}

	for _, testCase := range testCases {
		val := BCT(testCase.f())
		t.Logf("test case name: %s", testCase.name)
		require.Equal(t, testCase.expect, val)
	}
}

func newValidCompleteTree() *Node {
	root := &Node{Val: 1}
	root.Left = &Node{Val: 2}
	root.Right = &Node{Val: 3}
	root.Left.Left = &Node{Val: 4}
	return root
}

func newInvalidCompleteTree() *Node {
	root := &Node{Val: 1}
	root.Left = &Node{Val: 2}
	root.Right = &Node{Val: 3}
	root.Right.Left = &Node{Val: 7}
	root.Right.Right = &Node{Val: 8}

	return root
}
