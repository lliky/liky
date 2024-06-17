package binary_tree

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBCT(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		f      func() *TreeNode
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

func newValidCompleteTree() *TreeNode {
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Left.Left = &TreeNode{Val: 4}
	return root
}

func newInvalidCompleteTree() *TreeNode {
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Right.Left = &TreeNode{Val: 7}
	root.Right.Right = &TreeNode{Val: 8}

	return root
}
