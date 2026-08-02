package main

import (
	"fmt"
)

// Use slice when number of elements can change
// Slice is suitable for ordered data and accessing elements by index
// A slice contains a pointer to an underlying array, length, and capacity
func main() {

	// Slice does not have a fixed length
	numbers := []int{1, 2, 3}

	fmt.Println(numbers)      // [1 2 3]
	fmt.Println(len(numbers)) // 3
	fmt.Println(cap(numbers)) // 3

	// append can add more elements to a slice
	// It may create a new underlying array when capacity is not enough
	numbers = append(numbers, 4)

	fmt.Println(numbers)      // [1 2 3 4]
	fmt.Println(len(numbers)) // 4

	// Passing a slice to a function does not copy all elements
	// It copies only the slice header: pointer, length, and capacity
	// Both slices still refer to the same underlying array
	func(arr []int) {
		arr[1] = 20

		fmt.Println(arr) // [1 20 3 4]
	}(numbers)

	// The original slice is also modified
	fmt.Println(numbers) // [1 20 3 4]

	// However, changing the slice length inside a function
	// does not change the caller's slice header
	func(arr []int) {
		arr = append(arr, 5)

		fmt.Println(arr) // [1 20 3 4 5]
	}(numbers)

	fmt.Println(numbers) // [1 20 3 4]

	// Return the new slice if we want the caller to receive
	// the result of append
	numbers = func(arr []int) []int {
		arr = append(arr, 5)

		return arr
	}(numbers)

	fmt.Println(numbers) // [1 20 3 4 5]

	// A sliced slice usually shares the same underlying array
	part := numbers[1:4]

	fmt.Println(part) // [20 3 4]

	part[0] = 200

	// Modifying part also modifies numbers
	fmt.Println(part)    // [200 3 4]
	fmt.Println(numbers) // [1 200 3 4 5]

	// Use copy when we need an independent slice
	independent := make([]int, len(part))
	copy(independent, part)

	independent[0] = 999

	fmt.Println(independent) // [999 3 4]
	fmt.Println(part)        // [200 3 4]
	fmt.Println(numbers)     // [1 200 3 4 5]

	// Slice cannot be compared with another slice using ==
	another := []int{1, 200, 3, 4, 5}

	// fmt.Println(numbers == another) // compile error

	// We need to compare each element manually
	fmt.Println(equalSlices(numbers, another)) // true

	// nil slice is valid and can be appended to
	var empty []int

	fmt.Println(empty)        // []
	fmt.Println(empty == nil) // true
	fmt.Println(len(empty))   // 0

	empty = append(empty, 10)

	fmt.Println(empty) // [10]
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
