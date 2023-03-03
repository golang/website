// +build OMIT

package main

import "fmt"

const Warm = "to 175°C" // string constant
const Pie = false       // boolean constant
const Hash = '#'        // rune constant
const One = 1           // integer constant
const Pi = 3.14         // floating-point constant
const Other = 6.7i      // imaginary constant

func main() {
	const World = "世界"
	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")

	const Truth = true
	fmt.Println("Go rules?", Truth)
}
