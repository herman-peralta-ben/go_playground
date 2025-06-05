package main

import (
	"fmt"
)

// region A
type A struct{}

// Hello is a method defined on A.
func (A) Hello() {
	fmt.Println("Hello from A")
}

// Bye is another method defined on A.
func (A) Bye() {
	fmt.Println("Bye from A")
}

func (A) Internal() int {
	return 42
}

//endregion A

// region B
// B embeds A, inheriting its methods via composition.
type B struct {
	A
}

// Hello overrides A's Hello method in the context of B.
// This hides A.Hello when calling Hello on a B instance.
func (B) Hello() {
	fmt.Println("Hello from B")
}

// Call A.Internal by holding a "self/this" reference to B
// By convention 's' is used instead of 'self'.
func (self *B) Other() {
	fmt.Println(self.Internal())
}

//endregion B

// $ go run basics/structs/embedding/main.go
func main() {
	b := B{}
	b.Hello() // → Hello from B (B's override)
	b.Bye()   // → Bye from A (inherited from A)
	b.Other() // → 42 ()
}
