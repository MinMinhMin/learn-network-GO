package main

import (
	"fmt"
)

// Use array when number of elements is fixed and have it meaning
// Because it have fixed lenght -> fixed memory -> easy to handle/manage memory
func main() {

	// we can compare two array (which in map/slice we cannot)
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}

	fmt.Println(a == b) // true

	// lenght is also a feature of array, this function only accept array with length of 3
	func(arr [3]int) {
		arr[1] = 3

		// pass array to function would make a copy -> increase memory

		fmt.Println(arr) // [1 3 3]
	}(a)

	fmt.Println(a) // [1 2 3]

	// we can modify our actual array with pointer
	func(arr *[3]int) {
		arr[1] = 4

		fmt.Println(arr) // [1 4 3]
	}(&a)

	fmt.Println(a) // [1 4 3]

}
