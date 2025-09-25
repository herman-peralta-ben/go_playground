// go run hello_world.go

package main

import (
	"fmt" // serves as prefix for all functions of the package
)

/* This is a multiline comment */
//This is a single line comment

// Entry point, package main and main function
func main() {
	// This will be called when main finishes, even with error
	defer func() {
		fmt.Println("finished main")
	}()

	name := "World"
	fmt.Printf("Hello %s\n", name)
	// All the package exported names should start with a capital letter.

	dataStructures()
	loops()
	conditionals()
	switches()

	func() {
		fmt.Println("Annonymous Function / Lambda")
	}() // --> () means run this annonymous function

	lambda := func(name string) {
		fmt.Printf("Lambda: Hello %s\n", name)
	}

	lambda("Herman")

	goto gotolabel

gotolabel:
	fmt.Println("Inside goto statement")

	declarations()
}

func dataStructures() {
	fmt.Println("================= Data structures =================")
	var str = "hello world"
	const num = 42
	arr := [6]int{10, 20, 30, 40, 50, 60}
	map_data := map[string]int{
		"hello": 80,
		"world": 45,
	}
	fmt.Println(str, num, arr, map_data)
}

func loops() {
	fmt.Println("================= Loops =================")
	arr := [6]int{10, 20, 30, 40, 50, 60}
	// [range] is used with [for loop] to iterate over a
	// specific given data.
	for pos, val := range arr {
		if pos == 5 {
			break
		}
		if pos == 1 {
			continue
		}
		fmt.Printf("value at index %d is %d\n", pos, val)
	}
}

func conditionals() {
	fmt.Println("================= Conditionals =================")
	const num = 42
	if num > 100 {
		fmt.Println("num is greater than 100")
	} else if num == 42 {
		fmt.Println("num is the answer 🤓")
	} else {
		fmt.Println("num is not greater than 100")
	}
}

func switches() {
	fmt.Println("================= Switch =================")
	var num = 5
	switch num {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	default:
		fmt.Println("No match")
	}

	switch num {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 10:
		fmt.Println("three")
		fallthrough
		// [fallthrough] transfers the program control to the next case (or default case),
		// skipping the given condition.
	case 11:
		fmt.Println("next case")
	default:
		fmt.Println("default statement")
	}
}

func declarations() {
	var variable7, variable8, variable9 int
	variable7, variable8, variable9 = 123, 456, 789
	fmt.Println(variable7, variable8, variable9)

	str, num := "multiple", 2
	fmt.Println(str, num)

	var (
		variable10 int
		variable11 string
		variable12 bool
		_          int // _ ignores
	)
	fmt.Println(variable10, variable11, variable12)
}
