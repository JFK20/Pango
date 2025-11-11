// main.go
package main

import (
	"fmt"
	"pango/series" // replace with your module path
)

func main() {
	// Create a series: ages with custom index
	ages := series.NewSeries(
		"Age",                               // name
		[]any{25, 30, 35},                   // values
		[]string{"Alice", "Bob", "Charlie"}, // index

	)

	fmt.Println(ages)
	// Output:
	// Age
	// Alice: 25
	// Bob: 30
	// Charlie: 35

	fmt.Println(ages.Head(2))
	fmt.Println(ages.Tail(2))
}
