package main

import (
	"fmt"
	"sort"
)

// Use map when we need to store and find values by key
// Map is suitable for lookup, counting, grouping, and checking existence
// Reading, inserting, and deleting usually have average O(1) complexity
func main() {

	// Create and initialize a map
	ages := map[string]int{
		"An":   16,
		"Binh": 17,
	}

	fmt.Println(ages)         // map[An:16 Binh:17]
	fmt.Println(ages["An"])   // 16
	fmt.Println(ages["Binh"]) // 17

	// Add or update a value
	ages["Chi"] = 18
	ages["An"] = 20

	fmt.Println(ages["An"]) // 20

	// When a key does not exist, map returns the zero value
	fmt.Println(ages["Dung"]) // 0

	// Use the comma-ok syntax to check whether a key exists
	age, exists := ages["Dung"]

	fmt.Println(age)    // 0
	fmt.Println(exists) // false

	if age, exists := ages["An"]; exists {
		fmt.Println("An exists with age:", age)
	}

	// Delete an element
	delete(ages, "Binh")

	fmt.Println(ages)

	// Deleting a key that does not exist is safe
	delete(ages, "Unknown")

	// Passing a map to a function does not copy all map data
	// Changes made inside the function affect the original map
	func(values map[string]int) {
		values["An"] = 25
		values["David"] = 30
	}(ages)

	fmt.Println(ages["An"])    // 25
	fmt.Println(ages["David"]) // 30

	// Map cannot be compared with another map using ==
	another := map[string]int{
		"An":    25,
		"Chi":   18,
		"David": 30,
	}

	// fmt.Println(ages == another) // compile error

	// We need to compare keys and values manually
	fmt.Println(equalMaps(ages, another)) // true

	// Map iteration order is not guaranteed
	for key, value := range ages {
		fmt.Println(key, value)
	}

	// If we need stable order, copy keys into a slice and sort them
	keys := make([]string, 0, len(ages))

	for key := range ages {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Println(key, ages[key])
	}

	// Map can be used as a set
	// struct{} uses no extra memory for the value
	visited := map[string]struct{}{}

	visited["page-1"] = struct{}{}
	visited["page-2"] = struct{}{}

	_, pageExists := visited["page-1"]

	fmt.Println(pageExists) // true

	// Map is commonly used for counting
	words := []string{"go", "java", "go", "python", "go"}
	counts := make(map[string]int)

	for _, word := range words {
		counts[word]++
	}

	fmt.Println(counts["go"])     // 3
	fmt.Println(counts["java"])   // 1
	fmt.Println(counts["python"]) // 1

	// A nil map can be read safely
	var nilMap map[string]int

	fmt.Println(nilMap["test"]) // 0
	fmt.Println(nilMap == nil)  // true

	// But writing to a nil map causes panic
	// nilMap["test"] = 1

	// Initialize it before writing
	nilMap = make(map[string]int)
	nilMap["test"] = 1

	fmt.Println(nilMap["test"]) // 1

	// Struct values stored in a map cannot have fields modified directly
	users := map[string]User{
		"u1": {
			Name: "An",
			Age:  16,
		},
	}

	// users["u1"].Age = 17 // compile error

	// Get the value, modify it, and assign it back
	user := users["u1"]
	user.Age = 17
	users["u1"] = user

	fmt.Println(users["u1"]) // {An 17}

	// Another option is storing pointers in the map
	userPointers := map[string]*User{
		"u2": {
			Name: "Binh",
			Age:  18,
		},
	}

	userPointers["u2"].Age = 19

	fmt.Println(*userPointers["u2"]) // {Binh 19}
}

type User struct {
	Name string
	Age  int
}

func equalMaps(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valueA := range a {
		valueB, exists := b[key]

		if !exists || valueA != valueB {
			return false
		}
	}

	return true
}
