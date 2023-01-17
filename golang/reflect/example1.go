package reflect

import (
	"fmt"
	"reflect"
)

func example1() {
	var a int
	fmt.Println("type: ", reflect.TypeOf(a))
	fmt.Println("value: ", reflect.ValueOf(a))
}
