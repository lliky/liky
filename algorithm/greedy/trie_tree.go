package greedy

type TrieNode struct {
	pass int
	end  int
	// nexts[0] == nil 没有走向 'a' 的路
	// nexts[1] != nil 有走向 'a' 的路
	// ...
	// nexts[25] != nil 有走向 'z' 的路
	nexts []*TrieNode // map<char, TrieNode> nexts
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		nexts: make([]*TrieNode, 26),
	}
}

// NewTrie 头节点
func NewTrie() *TrieNode {
	return newTrieNode()
}

// Insert 加入单词
func (root *TrieNode) Insert(word string) {
	if word == "" {
		return
	}
	node := root
	node.pass++
	for i := 0; i < len(word); i++ { // 遍历字符
		index := word[i] - 'a' // 由字符，找哪一条路
		if node.nexts[index] == nil {
			node.nexts[index] = newTrieNode()
		}
		node = node.nexts[index]
		node.pass++
	}
	node.end++
}

// Search word 这个单词加入过几次
func (root *TrieNode) Search(word string) int {
	if len(word) == 0 {
		return 0
	}
	node := root
	for i := 0; i < len(word); i++ {
		index := word[i] - 'a'
		if node.nexts[index] == nil {
			return 0
		}
		node = node.nexts[index]
	}
	return node.end
}

// Delete 删除该单词，同一个单词需要删除多次
func (root *TrieNode) Delete(word string) {
	if root.Search(word) != 0 { // 表示存在，就删除
		node := root
		node.pass--
		for i := 0; i < len(word); i++ {
			index := word[i] - 'a'
			node.nexts[index].pass--
			if node.nexts[index].pass == 0 { // 说明不需要
				node.nexts[index] = nil
				return
			}
			node = node.nexts[index]
		}
		node.end--
	}
}

// PrefixNumber 加入所有字符串中，有几个是以 pre 这个字符串作为前缀的
func (root *TrieNode) PrefixNumber(pre string) int {
	if len(pre) == 0 {
		return 0
	}
	node := root
	for i := 0; i < len(pre); i++ {
		index := pre[i] - 'a'
		if node.nexts[index] == nil {
			return 0
		}
		node = node.nexts[index]
	}
	return node.pass
}
