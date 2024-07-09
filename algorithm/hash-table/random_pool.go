package hash_table

import "golang.org/x/exp/rand"

type Pool struct {
	keyIndexMap map[int]int
	indexKeyMap map[int]int
	size        int
}

func NewPool() Pool {
	return Pool{
		keyIndexMap: make(map[int]int),
		indexKeyMap: make(map[int]int),
	}
}

func (p *Pool) Insert(val int) {
	if _, ok := p.keyIndexMap[val]; !ok {
		p.keyIndexMap[val] = p.size
		p.indexKeyMap[p.size] = val
		p.size++
	}
}

func (p *Pool) Delete(val int) {
	if _, ok := p.keyIndexMap[val]; ok {
		deleteIndex := p.keyIndexMap[val]

		lastIndex := p.size - 1
		lastKey := p.indexKeyMap[lastIndex]

		p.keyIndexMap[lastKey] = deleteIndex
		p.indexKeyMap[deleteIndex] = lastKey

		delete(p.keyIndexMap, val)
		delete(p.indexKeyMap, lastIndex)
		p.size--
	}
}

func (p *Pool) GetRandom() int {
	random := rand.Intn(p.size)
	return p.indexKeyMap[random]
}
