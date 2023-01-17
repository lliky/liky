package reflect

import (
	"fmt"
	"reflect"
)

func example4() {
	var num = 1.2345
	fmt.Println("the old value of num : ", num)

	pointer := reflect.ValueOf(&num)
	newValue := pointer.Elem()

	fmt.Println("type of pointer: ", newValue.Type())
	fmt.Println("settability of pointer: ", newValue.CanSet())

	newValue.SetFloat(2.0)
	fmt.Println("the new value of num", num)

	pointer = reflect.ValueOf(num)
	// panic, 这里必须是指针
	// newValue = pointer.Elem()
}
