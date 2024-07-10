package kmp

import (
	"fmt"
	"testing"
)

func TestStrStr(t *testing.T) {
	s := "aaaaaabaaa"
	substr := "baaa"
	fmt.Println(StrStr(s, substr))
}
