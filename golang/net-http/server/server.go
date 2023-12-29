package main

import (
	"fmt"
	"net/http"
)

type DemoHandler struct {
}

func (d *DemoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "hello world")
}

func main() {
	http.Handle("/ping", &DemoHandler{})
	http.ListenAndServe(":8088", nil)
}
