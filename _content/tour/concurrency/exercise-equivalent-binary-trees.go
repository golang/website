// +build no-build OMIT

package main

import "golang.org/x/tour/tree"

// Walk 遍历树 t，并树中所有的值发送到信道 ch。
func Walk(t *tree.Tree, ch chan int)

// Same 判断 t1 和 t2 是否包含相同的值。
func Same(t1, t2 *tree.Tree) bool

func main() {
}
