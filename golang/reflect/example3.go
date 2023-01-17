package reflect

import (
	"fmt"
	"reflect"
)

type User struct {
	Id   int
	Age  int
	Name string
}

func (u User) ReflectCallFunc() {
	fmt.Println("hello world")
}
func example3() {
	user := User{1, 2, "hello"}
	DoFieldAndMethod(user)
}

func DoFieldAndMethod(in interface{}) {
	getType := reflect.TypeOf(in)
	fmt.Println("get Type is: ", getType.Name())

	getValue := reflect.ValueOf(in)
	fmt.Println("get Value is: ", getValue)
	// 获取字段
	// 1. 先获取 interface 的 reflect.Type, 然后通过 NumField 进行遍历
	// 2. 再通过 reflect.Type 的 Field 获取其 Field
	// 3. 再通过 reflect.Value 的 Field 的 interface() 得到对应的 value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}
	// 获取字段
	// 先获取 interface 的 reflect.Type，然后通过 NumMethod 进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s : %v\n", m.Name, m.Type)
	}
}
