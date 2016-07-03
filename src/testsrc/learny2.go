// from https://learnxinyminutes.com/docs/go/
// Single line comment
/* Multi-
 line comment */

// A package clause starts every source file.
// Main is a special name declaring an executable rather than a library.
package main

// Import declaration declares library packages referenced in this file.
import (
    "fmt"       // A package in the Go standard library.
    "io/ioutil" // Implements some I/O utility functions.
    m "math"    // Math library with local alias m.
    "net/http"  // Yes, a web server!
    "strconv"   // String conversions.
)

// A function definition. Main is special. It is the entry point for the
// executable program. Love it or hate it, Go uses brace brackets.
func main() {
    // Println outputs a line to stdout.
    // Qualify it with the package name, fmt.
    fmt.Println("Hello world!")

    // Call another function within this package.
    beyondHello()
}

// Functions have parameters in parentheses.
// If there are no parameters, empty parentheses are still required.
func beyondHello() {
    var x int // Variable declaration. Variables must be declared before use.
    x = 3     // Variable assignment.
	var (
		v1, v2 int = 4, 5
		v3 = 5,
		v4
	)

    // "Short" declarations use := to infer the type, declare, and assign.
    y := 4
    y, prod := learnMultiple(x, y)        // Function returns two values.
    fmt.Println("sum:", sum, "prod:", prod) // Simple output.
    learnTypes()                            // < y minutes, learn more!
}
