package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	baseString := "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(baseString))
}
