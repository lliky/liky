package reflect

import (
	"fmt"
	"reflect"
)

type User1 struct {
	Id   int
	Name string
	Age  int
}

func (u User1) ReflectCallFuncHasArgs(name string, age int) {
	fmt.Printf("Get args name: %v, age: %v; and original Name: %v\n", name, age, u.Name)
}

func (u User1) ReflectCallFuncNoArgs() {
	fmt.Println("ReflectCallFuncNoArgs")
}

func example5() {
	user := User1{1, "hi", 2}
	getValue := reflect.ValueOf(user)
	// 通过 MethodByName 进行注册
	methodValue := getValue.MethodByName("ReflectCallFuncHasArgs")
	args := []reflect.Value{reflect.ValueOf("he"), reflect.ValueOf(33)}
	methodValue.Call(args)
	methodValue = getValue.MethodByName("ReflectCallFuncNoArgs")
	args = make([]reflect.Value, 0)
	methodValue.Call(args)
	fmt.Println(user)
}
