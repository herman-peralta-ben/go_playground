package main

import (
	"fmt"
	"math"
	"unsafe"
)

func main() {
	fmt.Println("************************ Primitive types ************************")
	primitiveDataTypes()
	fmt.Println("********************** Non Primitive types **********************")
	nonPrimitiveDataTypes()
}

func primitiveDataTypes() {
	numericDataTypes()
	runeDataType()
	stringDataType()
}

// aka Derived or Referenced data types (contain the address
// of the dynamically created objects in memory.)
func nonPrimitiveDataTypes() {
	pointerDataType()
	arrayDataType()
	structDataType()
	sliceDataType()
}

func numericDataTypes() {
	fmt.Println("================ Numeric Data Types ================")

	//# uint
	lowUInt8, highUInt8 := uint8(0), uint8(math.MaxUint8)
	fmt.Printf("uint8: [%d, %d]\n", lowUInt8, highUInt8)
	lowUInt16, highUInt16 := uint16(0), uint16(math.MaxUint16)
	fmt.Printf("uint16: [%d, %d]\n", lowUInt16, highUInt16)
	lowUInt32, highUInt32 := uint32(0), uint32(math.MaxUint32)
	fmt.Printf("uint32: [%d, %d]\n", lowUInt32, highUInt32)
	lowUInt64, highUInt64 := uint64(0), uint64(math.MaxUint64)
	fmt.Printf("uint64: [%d, %d]\n", lowUInt64, highUInt64)

	//# int
	lowInt8, highInt8 := int8(math.MinInt8), int8(math.MaxInt8)
	fmt.Printf("int8: [%d, %d]\n", lowInt8, highInt8)
	lowInt16, highInt16 := int16(math.MinInt16), int16(math.MaxInt16)
	fmt.Printf("int16: [%d, %d]\n", lowInt16, highInt16)
	lowInt32, highInt32 := int32(math.MinInt32), int32(math.MaxInt32)
	fmt.Printf("int32: [%d, %d]\n", lowInt32, highInt32)
	lowInt64, highInt64 := int64(math.MinInt64), int64(math.MaxInt64)
	fmt.Printf("int64: [%d, %d]\n", lowInt64, highInt64)

	//# Float
	lowFloat32, highFloat32 := float32(math.SmallestNonzeroFloat32), float32(math.MaxFloat32)
	fmt.Printf("float32: [%f, %f]\n", lowFloat32, highFloat32)
	lowFloat64, highFloat64 := float64(math.SmallestNonzeroFloat64), float64(math.MaxFloat64)
	fmt.Printf("float64: [%f, %f]\n", lowFloat64, highFloat64)

	//# complex
	minComplex64, maxComplex64 := complex64(complex(float32(math.SmallestNonzeroFloat32), float32(math.SmallestNonzeroFloat32))), complex64(complex(float32(math.MaxFloat32), float32(math.MaxFloat32)))
	fmt.Printf("complex64: [%v, %v]\n", minComplex64, maxComplex64)
	minComplex128, maxComplex128 := complex128(complex(float64(math.SmallestNonzeroFloat64), float64(math.SmallestNonzeroFloat64))), complex128(complex(float64(math.MaxFloat64), float64(math.MaxFloat64)))
	fmt.Printf("complex128: [%v, %v]\n", minComplex128, maxComplex128)
}

func runeDataType() {
	fmt.Println("=================== Rune =====================")
	//# rune, is an alias dor int32 used to identify a unique Unicode character.
	// under the hood rune is a number
	var aRune rune = 'a'
	fmt.Println(aRune)        // Prints 97 (Unicode)
	fmt.Printf("%c\n", aRune) // Prints a
	var kaRune rune = 'か'
	fmt.Println(kaRune)        // Prints 12363 (Unicode)
	fmt.Printf("%c\n", kaRune) // Prints か
}

func stringDataType() {
	fmt.Println("=================== String =====================")
	//# string
	// In Go, a string is composed by UTF-8 bytes, not "chars". To process real characters (not bytes), use [rune].
	s := "mañana"
	for _, r := range s {
		fmt.Printf("Rune: %c, Code point: %U\n", r, r)
	}
	// 1. string → []byte
	bytes := []byte(s)
	fmt.Println(bytes) // UTF-8 bytes: [109 97 195 177 97 110 97], ñ requires two bytes [195 177]
	// 2. string → []rune
	runes := []rune(s)
	fmt.Println(runes)        // [109 97 241 97 110 97]
	fmt.Printf("%c\n", runes) // [m a ñ a n a]
	// 3. []byte → string
	byteArr := []byte{109, 97, 195, 177, 97, 110, 97}
	fmt.Println(string(byteArr)) // mañana
	// 4. []rune → string
	runeArr := []rune{109, 97, 241, 97, 110, 97}
	fmt.Println(string(runeArr)) // mañana
}

func pointerDataType() {
	fmt.Println("=================== Pointer =====================")

	// A pointer stores the address of a value in memory.
	// Its *type* depends on the type of the value it points to.
	// Example: *int, *string, *float64...

	// 🔎 The size of a pointer (the address) is always the same,
	// regardless of what it points to:
	// - 4 bytes on a 32-bit system
	// - 8 bytes on a 64-bit system

	var number int = 42
	fmt.Println("number:", number) // 42

	var pointerToNumber *int = &number

	fmt.Println("address of number:   ", &number)          // address
	fmt.Println("pointerToNumber:     ", pointerToNumber)  // same address
	fmt.Println("value via pointer:   ", *pointerToNumber) // 42

	*pointerToNumber += 10                                 // modify original value through the pointer
	fmt.Println("modified via pointer:", *pointerToNumber) // 52
	fmt.Println("number now:          ", number)           // 52

	// Iterating over an array using pointer arithmetic (like in C)
	/*
		C equivalent:
					```c
					int arr[] = {1, 2, 3};
					int *base = arr;
					for (int i = 0; i < 3; i++) {
					    printf("%d\n", *(base + i));
					}
					```
		In Go:
			- unsafe.Pointer: allows working with untyped Pointers
			- uintptr: allows pointer arithmetic, e.g. (base + 1*sizeof)
			- unsafe.Sizeof: gives the size in bytes of each array element
	*/
	func() {
		arr := [3]int{10, 20, 30}
		base := unsafe.Pointer(&arr[0])

		for i := 0; i < len(arr); i++ {
			// Cast the untyped Pointer back to int Pointer, so we can deference ut
			elementPtr := (*int)(unsafe.Pointer(uintptr(base) + uintptr(i)*unsafe.Sizeof(arr[0])))
			// Use %p for printing pointers
			fmt.Printf("By pointer: arr[%d]=%d (@ %p, %p)\n", i, *elementPtr, elementPtr, &arr[i])
		}
	}()
}

func arrayDataType() {
	fmt.Println("=================== Array =====================")

	// Declare an array of 5 integers (fixed length).
	var intArr = [5]int{0, 1, 2, 3, 4}

	// Print entire array and first element
	fmt.Println("Original array:", intArr)
	fmt.Println("First element :", intArr[0]) // 0

	// Arrays in Go are value types:
	// Assigning to another variable copies all elements
	copy := intArr
	fmt.Println("Copied array  :", copy)

	// Modify original array
	intArr[0] = 42
	fmt.Println("Modified original array:", intArr) // [42 1 2 3 4]
	fmt.Println("Unchanged copy         :", copy)   // [0 1 2 3 4]

	// Print the address of the original array
	fmt.Printf("Original intArr address      : %p\n", &intArr)

	// Pass array to function: it will be copied
	func(arr [5]int) {
		// The address inside is different => arr is a copy
		fmt.Printf("Function param arr address   : %p\n", &arr)
		fmt.Printf("Still printing outer intArr  : %p (same as before)\n", &intArr)

		// Change inside the copy
		arr[0] = 42
		fmt.Println("Modified inside func (copy)  :", arr)
	}(intArr)

	// Original remains unchanged
	fmt.Println("Original after func call     :", intArr) // still [42 1 2 3 4]

	// Pass by reference
	func(arr *[5]int) {
		arr[0] = 666
	}(&intArr)
	fmt.Println("Original after reference func call:", intArr) // updated [666 1 2 3 4]
}

func structDataType() {
	// private: start with lowercase
	// Public: start with uppercase
	type PublicStruct struct {
		id   string // private
		Name string // public
	}

	var str PublicStruct // Create
	str.Name = "Gopher"
	str.id = "10"

	fmt.Println(str)
}

// Slices are growable lists
// Passed by reference by default on functions
func sliceDataType() {
	slice1 := make([]int, 0)
	fmt.Println(slice1)
	slice2 := make([]int, 3)
	slice2[0] = 0
	slice2[1] = 1
	slice2[2] = 2
	fmt.Println(slice2) // [0 1 2]
	slice3 := []string{"abc", "def", "ghi"}
	fmt.Println(slice3) //  [abc def ghi]
	slice3 = append(slice3, "jkl")
	fmt.Println(slice3) // [abc def ghi jkl]
	fmt.Println(&slice3)
	slice4 := make([]string, 5)
	copy(slice4, slice3)
	fmt.Println(slice4) // [abc def ghi jkl]
	fmt.Println(&slice4)
	slice2D := make([][]string, 2, 3)
	// len = 2 → only slice2D[0] and slice2D[1] can be accessed
	// cap = 3 → grow up to 3 elements, (but need to use **append**)
	fmt.Println(slice2D) //  [[] []]
	slice2D[0] = []string{"a", "b", "c"}
	slice2D[1] = []string{"d", "e"}
	fmt.Println(slice2D) //  [[a b c] [d e]]
	// slice2D[2] = []string{"f", "g"} ---> Error, need to use append()
	slice2D = append(slice2D, []string{"f", "g"})
	fmt.Println(slice2D) //  [[a b c] [d e] [f g]]
}
