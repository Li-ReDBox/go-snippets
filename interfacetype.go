package main

import "fmt"

type Base struct {
	X, Y int
}

func (t Base) display() {
	fmt.Printf("%T display method\n", t)
	fmt.Println(t.X, t.Y)
}

// embedding the type which has the commone fields and methods
type SplitA struct{ Base }
type SplitB struct{ Base }

// the abstraction of common method but with different implementation
type Executer interface {
	commonDisplay()
}

// Here each different concrete type implement the interface
// to uniform the difference between them by wrapping
func (sa SplitA) commonDisplay() {
	fmt.Println("This is split .. Aa")
	sa.display()
}

func (sa SplitB) commonDisplay() {
	fmt.Println("This is split ... Bb")
	sa.display()
}

// we have to use a third party function - basically a helper
// we cannot implement this as a function back to any concret
// types.
func execute(t Executer) {
	fmt.Println("This is the helper function executing interface")
	fmt.Printf("In helper, the type it sees: %T\n", t)
	t.commonDisplay()
}

func main() {
	a := SplitA{Base{1, 2}}
	b := SplitB{Base{10, 20}}
	execute(a)
	fmt.Println()
	execute(b)
}
