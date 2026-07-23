// main.go
package main

import (
	"fmt"

	"pango/example/dataframes"
	"pango/example/series"
)

func main() {
	fmt.Println("=== Series examples ===")

	fmt.Println("\n--- Example 1: Column sum ---")
	series.Example1SumColumn()

	fmt.Println("\n--- Example 2: add two columns and save it in a third ---")
	series.Example2AddColumns()

	fmt.Println("\n--- Example 3: multiple arithmetic operations ---")
	series.Example3MultipleOperations()

	fmt.Println("\n=== DataFrame examples ===")

	fmt.Println("\n--- Example 1: Select, filter, Head/Tail ---")
	dataframes.Example1SelectFilterHeadTail()

	fmt.Println("\n--- Example 2: Transform columns ---")
	dataframes.Example2TransformColumns()

	fmt.Println("\n--- Example 3: GroupBy and aggregate ---")
	dataframes.Example3GroupByAggregate()
}
