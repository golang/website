// +build OMIT

package main

import "fmt"

func helper() string {
	fmt.Println("helper")
	return "world"
}

func main() {
	defer fmt.Println(helper())

	fmt.Println("hello")
}
