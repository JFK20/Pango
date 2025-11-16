// main.go
package main

import (
	"fmt"
	"pango/Dataframe"
	"pango/series"
)

func main() {
	fmt.Println("=== Example 1: Coloum sum ===")
	example1SumColumn()

	fmt.Println("\n=== Example 2: add to columns and save it in a third ===")
	example2_AddColumns()

	fmt.Println("\n=== Example 3: multiple arithmetic operations ===")
	example3_MultipleOperations()
}

// Beispiel 1: Eine Spalte summieren
func example1SumColumn() {
	// Erstelle Series für ein DataFrame
	names := series.NewIndexSeries("Name", []string{"Alice", "Bob", "Charlie", "David"})
	sales := series.NewIndexNumericSeries("sales", []int{100, 150, 200, 175})
	revenue := series.NewIndexNumericSeries("revenue", []float64{1500.50, 2250.75, 3000.00, 2625.25})

	// Erstelle DataFrame
	df, err := dataframe.NewDataFrame(names, sales, revenue)
	if err != nil {
		fmt.Println("Error during DataFrame creation:", err)
		return
	}

	fmt.Println("DataFrame:")
	fmt.Println(df)

	// sum the sales-column
	salesCol, err := df.GetNumericColumn("sales")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	salesSum := salesCol.SumFloat()
	fmt.Printf("\nsumSales: %.0f\n", salesSum)

	// sum the revenue value
	revenueCol, err := df.GetNumericColumn("revenue")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	revenueSum := revenueCol.SumFloat()
	fmt.Printf("Total Revenue: %.2f €\n", revenueSum)
}

// Beispiel 2: Zwei Spalten addieren und in einer dritten speichern
func example2_AddColumns() {
	// Erstelle Mitarbeiter-DataFrame mit Grundgehalt und Bonus
	mitarbeiter := series.NewIndexSeries("Mitarbeiter", []string{"Anna", "Ben", "Clara", "Daniel"})
	grundgehalt := series.NewIndexNumericSeries("Grundgehalt", []float64{3000, 3500, 4000, 3200})
	bonus := series.NewIndexNumericSeries("Bonus", []float64{500, 750, 1000, 600})

	// Erstelle DataFrame
	df, err := dataframe.NewDataFrame(mitarbeiter, grundgehalt, bonus)
	if err != nil {
		fmt.Println("Fehler:", err)
		return
	}

	fmt.Println("DataFrame (vorher):")
	fmt.Println(df)

	// Hole die beiden Spalten
	grundgehaltCol, _ := df.GetNumericColumn("Grundgehalt")
	bonusCol, _ := df.GetNumericColumn("Bonus")

	// Konvertiere zu float64 NumericSeries
	grundgehaltFloat := grundgehaltCol.(*series.NumericSeries[float64, int])
	bonusFloat := bonusCol.(*series.NumericSeries[float64, int])

	// Addiere die beiden Spalten
	gesamtgehalt := grundgehaltFloat.Add(bonusFloat, "Gesamtgehalt")

	// Füge die neue Spalte zum DataFrame hinzu
	err = df.AddColumn(gesamtgehalt)
	if err != nil {
		fmt.Println("Fehler beim Hinzufügen der Spalte:", err)
		return
	}

	fmt.Println("\nDataFrame (nachher - mit Gesamtgehalt):")
	fmt.Println(df)

	// Zeige die Summe aller Gesamtgehälter
	gesamtgehaltCol, _ := df.GetNumericColumn("Gesamtgehalt")
	totalSum := gesamtgehaltCol.SumFloat()
	fmt.Printf("\nGesamte Lohnkosten: %.2f €\n", totalSum)
}

// Example 3: Multiple arithmetic operations
func example3_MultipleOperations() {
	// Create product DataFrame
	products := series.NewIndexSeries[string]("Product", []string{"Widget A", "Widget B", "Widget C"})
	price := series.NewIndexNumericSeries("Price", []float64{19.99, 29.99, 39.99})
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
