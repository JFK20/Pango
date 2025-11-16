// main.go
package main

import (
	"fmt"
	"pango/Dataframe"
	"pango/series"
)

func main() {
	fmt.Println("=== Example 1: Column sum ===")
	example1SumColumn()

	fmt.Println("\n=== Example 2: add two columns and save it in a third ===")
	example2AddColumns()

	fmt.Println("\n=== Example 3: multiple arithmetic operations ===")
	example3Multipleoperations()
}

// Example 1: Sum a column
func example1SumColumn() {
	// Create Series for a DataFrame
	names := series.NewIndexSeries("Name", []string{"Alice", "Bob", "Charlie", "David"})
	sales := series.NewIndexNumericSeries("sales", []int{100, 150, 200, 175})
	revenue := series.NewIndexNumericSeries("revenue", []float64{1500.50, 2250.75, 3000.00, 2625.25})

	// Create DataFrame
	df, err := dataframe.NewDataFrame(names, sales, revenue)
	if err != nil {
		fmt.Println("Error during DataFrame creation:", err)
		return
	}

	fmt.Println("DataFrame:")
	fmt.Println(df)

	// sum the sales column
	salesCol, err := df.GetNumericColumn("sales")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	salesSum := salesCol.SumFloat()
	fmt.Printf("\nTotal Sales: %.0f\n", salesSum)

	// sum the revenue value
	revenueCol, err := df.GetNumericColumn("revenue")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	revenueSum := revenueCol.SumFloat()
	fmt.Printf("Total Revenue: %.2f €\n", revenueSum)
}

// Example 2: Add two columns and save in a third
func example2AddColumns() {
	// Create employee DataFrame with base salary and bonus
	employees := series.NewIndexSeries("employees", []string{"Anna", "Ben", "Clara", "Daniel"})
	baseSalary := series.NewIndexNumericSeries("baseSalary", []float64{3000, 3500, 4000, 3200})
	bonus := series.NewIndexNumericSeries("Bonus", []float64{500, 750, 1000, 600})

	// Create DataFrame
	df, err := dataframe.NewDataFrame(employees, baseSalary, bonus)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("DataFrame (before):")
	fmt.Println(df)

	// Get the two columns
	baseWageCol, _ := df.GetNumericColumn("baseSalary")
	bonusCol, _ := df.GetNumericColumn("Bonus")

	// Convert to float64 NumericSeries
	baseWageFloat := baseWageCol.(*series.NumericSeries[float64, int])
	bonusFloat := bonusCol.(*series.NumericSeries[float64, int])

	// Add the two columns
	totalSalary := baseWageFloat.Add(bonusFloat, "TotalSalary")

	// Add the new column to the DataFrame
	err = df.AddColumn(totalSalary)
	if err != nil {
		fmt.Println("Error adding column:", err)
		return
	}

	fmt.Println("\nDataFrame (after - with TotalSalary):")
	fmt.Println(df)

	// Show the sum of all total salaries
	totalSalaryCol, _ := df.GetNumericColumn("TotalSalary")
	totalSum := totalSalaryCol.SumFloat()
	fmt.Printf("\nTotal Labor Costs: %.2f €\n", totalSum)
}

// Example 3: Multiple arithmetic operations
func example3Multipleoperations() {
	// Create product DataFrame
	products := series.NewIndexSeries[string]("Product", []string{"Widget A", "Widget B", "Widget C"})
	price := series.NewIndexNumericSeries("Price", []float64{20, 30.50, 38.09})
	quantity := series.NewIndexNumericSeries("Quantity", []int{100, 150, 75})

	df, err := dataframe.NewDataFrame(products, price, quantity)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("DataFrame (Products):")
	fmt.Println(df)

	// Get the columns
	priceCol, _ := df.GetNumericColumn("Price")
	quantityCol, _ := df.GetNumericColumn("Quantity")

	// Convert to appropriate types and calculate total value
	// Note: For multiplication we need to align the types
	priceFloat := priceCol.(*series.NumericSeries[float64, int])
	quantityInt := quantityCol.(*series.NumericSeries[int, int])

	// Convert quantity to float64 for calculation
	quantityValues := make([]float64, quantityInt.Len())
	for i := 0; i < quantityInt.Len(); i++ {
		quantityValues[i] = float64(quantityInt.Values()[i])
	}
	quantityFloat := series.NewIndexNumericSeries("Quantity_float", quantityValues)

	// Calculate total value (Price * Quantity)
	totalValue := priceFloat.Multiply(quantityFloat, "TotalValue")

	// Add to DataFrame
	err = df.AddColumn(totalValue)
	if err != nil {
		fmt.Println("Error adding column:", err)
		return
	}

	fmt.Println("\nDataFrame (with TotalValue):")
	fmt.Println(df)

	// Sum the total value
	totalValueCol, _ := df.GetNumericColumn("TotalValue")
	sumTotalValue := totalValueCol.SumFloat()
	fmt.Printf("\nTotal Value of All Products: %.2f €\n", sumTotalValue)

	// Calculate average price
	avgPrice := priceFloat.Mean()
	fmt.Printf("Average Price: %.2f €\n", avgPrice)

	// Show Min and Max values
	fmt.Printf("Minimum Price: %.2f €\n", priceFloat.Min())
	fmt.Printf("Maximum Price: %.2f €\n", priceFloat.Max())
}
