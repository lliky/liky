package manacher

import (
	"fmt"
	"testing"
)

func TestManacher(t *testing.T) {
	s1 := "1221"
	fmt.Println(string(manacherString(s1)))
	s2 := "121"
	fmt.Println(string(manacherString(s2)))
}
