package reflect

import (
	"fmt"
	"reflect"
)

func example2() {
	var a = 2
	pointer := reflect.ValueOf(&a)
	value := reflect.ValueOf(a)

	convertPointer := pointer.Interface().(*int)
	convertValue := value.Interface().(int)
	// painc due to incomplete type
	// convertValue := value.Interface().(int)
	fmt.Printf("%v\n%v\n", convertPointer, convertValue)

	fmt.Printf("The type of pointer: %T\n", pointer)
	fmt.Printf("The type of value: %T\n", value)
	fmt.Printf("The type of convertPointer: %T\n", convertPointer)
	fmt.Printf("The type of convertVale: %T\n", convertValue)

}
