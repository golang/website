//go:build ignore || OMIT
// +build ignore OMIT

package main

import "golang.org/x/tour/reader"

type MyReader struct{}

// TODO: 为 MyReader 添加一个 Read([]byte) (int, error) 方法。

func main() {
	reader.Validate(MyReader{})
}
