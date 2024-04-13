//go:build OMIT
// +build OMIT

package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["答案"] = 42
	fmt.Println("值：", m["答案"])

	m["答案"] = 48
	fmt.Println("值：", m["答案"])

	delete(m, "答案")
	fmt.Println("值：", m["答案"])

	v, ok := m["答案"]
	fmt.Println("值：", v, "是否存在？", ok)
}
