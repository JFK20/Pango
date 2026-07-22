// Package dataframes contains runnable examples demonstrating
// pango/Dataframe-specific capabilities: selecting, filtering, slicing,
// transforming and grouping/aggregating DataFrames.
package dataframes

import (
	"fmt"

	dataframe "pango/Dataframe"
	"pango/series"
)

// Example1SelectFilterHeadTail demonstrates selecting a subset of columns,
// filtering rows by a predicate, and slicing with Head/Tail.
func Example1SelectFilterHeadTail() {
	products := series.NewIndexSeries("Product", []string{"Widget A", "Widget B", "Widget C", "Widget D", "Widget E"})
	category := series.NewIndexSeries("Category", []string{"Tools", "Tools", "Electronics", "Electronics", "Home"})
	price := series.NewIndexNumericSeries("Price", []float64{19.99, 24.99, 149.99, 89.99, 12.50})
	quantity := series.NewIndexNumericSeries("Quantity", []int{100, 80, 30, 45, 200})

	df, err := dataframe.NewDataFrame(products, category, price, quantity)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Full DataFrame:")
	fmt.Println(df)

	// Select only the columns we care about
	subset, err := df.Select("Product", "Price")
	if err != nil {
		fmt.Println("Error selecting columns:", err)
		return
	}
	fmt.Println("Selected columns (Product, Price):")
	fmt.Println(subset)

	// Filter for products priced above 50
	expensive, err := df.FilterColumn("Price", func(v any) bool {
		return v.(float64) > 50
	})
	if err != nil {
		fmt.Println("Error filtering:", err)
		return
	}
	fmt.Println("Products priced above 50:")
	fmt.Println(expensive)

	fmt.Println("First 2 rows (Head):")
	fmt.Println(df.Head(2))

	fmt.Println("Last 2 rows (Tail):")
	fmt.Println(df.Tail(2))
}

// Example2TransformColumns demonstrates ApplyToColumn, ApplyToColumns,
// RenameColumn and DropColumn.
func Example2TransformColumns() {
	products := series.NewIndexSeries("Product", []string{"Widget A", "Widget B", "Widget C"})
	price := series.NewIndexNumericSeries("Price", []float64{20, 30.5, 38.09})
	quantity := series.NewIndexNumericSeries("Quantity", []int{100, 150, 75})

	df, err := dataframe.NewDataFrame(products, price, quantity)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Apply a 10% discount to the Price column
	discounted, err := df.ApplyToColumn("Price", func(v any) any {
		return v.(float64) * 0.9
	})
	if err != nil {
		fmt.Println("Error applying discount:", err)
		return
	}
	fmt.Println("DataFrame with discounted prices:")
	fmt.Println(discounted)

	// Compute a total value column from Price and Quantity together
	withTotal, err := discounted.ApplyToColumns([]string{"Price", "Quantity"}, func(row map[string]any) any {
		return row["Price"].(float64) * float64(row["Quantity"].(int))
	}, "TotalValue")
	if err != nil {
		fmt.Println("Error computing total value:", err)
		return
	}
	fmt.Println("DataFrame with TotalValue column:")
	fmt.Println(withTotal)

	// Rename and drop columns
	if err := withTotal.RenameColumn("Product", "Item"); err != nil {
		fmt.Println("Error renaming column:", err)
		return
	}
	final, err := withTotal.DropColumn("Quantity")
	if err != nil {
		fmt.Println("Error dropping column:", err)
		return
	}
	fmt.Println("Final DataFrame (renamed + dropped column):")
	fmt.Println(final)
}

// Example3GroupByAggregate demonstrates grouping rows by a column and
// computing per-group aggregations and counts.
func Example3GroupByAggregate() {
	category := series.NewIndexSeries("Category", []string{"Tools", "Tools", "Electronics", "Electronics", "Home", "Home"})
	product := series.NewIndexSeries("Product", []string{"Widget A", "Widget B", "Gadget A", "Gadget B", "Lamp", "Rug"})
	revenue := series.NewIndexNumericSeries("Revenue", []float64{1200, 800, 5000, 3000, 450, 900})

	df, err := dataframe.NewDataFrame(category, product, revenue)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Sales DataFrame:")
	fmt.Println(df)

	gb, err := df.GroupBy("Category")
	if err != nil {
		fmt.Println("Error grouping:", err)
		return
	}

	agg, err := gb.AggregateColumns(map[string]map[string]dataframe.AggFunc{
		"Revenue": {
			"sum":  dataframe.AggSum,
			"mean": dataframe.AggMean,
			"max":  dataframe.AggMax,
		},
	})
	if err != nil {
		fmt.Println("Error aggregating:", err)
		return
	}
	fmt.Println("Revenue sum/mean/max by category:")
	fmt.Println(agg)

	counts, err := gb.Count()
	if err != nil {
		fmt.Println("Error counting groups:", err)
		return
	}
	fmt.Println("Number of products per category:")
	fmt.Println(counts)
}
