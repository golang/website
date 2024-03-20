// +build OMIT

package main

import "fmt"

func split(number int) (ones_digit, high_digits int) {
	// both named return values are already defined as int variables
	ones_digit = number % 10
	high_digits = number - ones_digit
	// will use named return values (ones_digit, high_digits)
	return
}

func main() {
	fmt.Println(split(17))
}
