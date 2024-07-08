package union

type Element struct {
	value interface{}
}

func NewElement(value interface{}) *Element {
	return &Element{value: value}
}

type Set struct {
	elementMap map[interface{}]*Element
	fatherMap  map[*Element]*Element
	sizeMap    map[*Element]int
}

func NewSet(list []interface{}) *Set {
	elementMap := make(map[interface{}]*Element)
	fatherMap := make(map[*Element]*Element)
	sizeMap := make(map[*Element]int)
	for _, v := range list {
		element := NewElement(v)
		elementMap[v] = element
		fatherMap[element] = element
		sizeMap[element] = 1
	}
	return &Set{elementMap: elementMap, fatherMap: fatherMap, sizeMap: sizeMap}
}

// FindHead 给定一个元素 element, 网上找，把代表元素返回
func (s *Set) FindHead(element *Element) *Element {
	path := make([]*Element, 0)
	for element != s.fatherMap[element] {
		path = append(path, element)
		element = s.fatherMap[element]
	}
	for _, v := range path {
		s.fatherMap[v] = element
	}
	return element
}

func (s *Set) IsSameSet(a, b interface{}) bool {
	_, oka := s.elementMap[a]
	_, okb := s.elementMap[b]
	if oka && okb {
		return s.FindHead(s.elementMap[a]) == s.FindHead(s.elementMap[b])
	}
	return false
}

func (s *Set) Union(a, b interface{}) {
	eA, oka := s.elementMap[a]
	eB, okb := s.elementMap[b]
	if oka && okb {
		aF := s.FindHead(eA)
		bF := s.FindHead(eB)
		big, small := aF, bF
		if s.sizeMap[aF] < s.sizeMap[bF] {
			big, small = bF, aF
		}
		s.fatherMap[small] = big
		s.sizeMap[big] = s.sizeMap[aF] + s.sizeMap[bF]
		delete(s.sizeMap, small)
	}
}
