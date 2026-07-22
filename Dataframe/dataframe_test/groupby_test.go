package dataframe_test

import (
	"math"
	"testing"

	dataframe "pango/Dataframe"
	"pango/series"
)

func newGroupedDataFrame(t *testing.T) *dataframe.DataFrame {
	t.Helper()

	category := series.NewIndexSeries("category", []string{"A", "A", "B", "B", "B"})
	product := series.NewIndexSeries("product", []string{"Widget", "Gadget", "Foo", "Bar", "Baz"})
	price := series.NewIndexNumericSeries("price", []float64{10, 20, 30, 40, 50})

	df, err := dataframe.NewDataFrame(category, product, price)
	if err != nil {
		t.Fatalf("unexpected error building grouped DataFrame: %v", err)
	}
	return df
}

func TestGroupBy(t *testing.T) {
	df := newGroupedDataFrame(t)

	if _, err := df.GroupBy("category"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.GroupBy("missing"); err == nil {
			t.Fatal("expected error grouping by a missing column")
		}
	})
}

func TestAggregateColumns(t *testing.T) {
	df := newGroupedDataFrame(t)

	gb, err := df.GroupBy("category")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := gb.AggregateColumns(map[string]map[string]dataframe.AggFunc{
		"price": {
			"sum":    dataframe.AggSum,
			"mean":   dataframe.AggMean,
			"min":    dataframe.AggMin,
			"max":    dataframe.AggMax,
			"stddev": dataframe.AggStdDev,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rows, cols := result.Shape(); rows != 2 || cols != 6 {
		t.Fatalf("expected shape (2, 6), got (%d, %d)", rows, cols)
	}

	// Regression: aggregation results used to always be nil because the
	// per-group sub-series lost numeric capability. Verify the real values.
	wantSum := map[string]float64{"A": 30, "B": 120}
	wantMean := map[string]float64{"A": 15, "B": 40}
	wantMin := map[string]float64{"A": 10, "B": 30}
	wantMax := map[string]float64{"A": 20, "B": 50}

	catCol, err := result.GetColumn("category")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sumCol, err := result.GetNumericColumn("price_sum")
	if err != nil {
		t.Fatalf("expected price_sum to be numeric: %v", err)
	}
	meanCol, err := result.GetNumericColumn("price_mean")
	if err != nil {
		t.Fatalf("expected price_mean to be numeric: %v", err)
	}
	minCol, err := result.GetNumericColumn("price_min")
	if err != nil {
		t.Fatalf("expected price_min to be numeric: %v", err)
	}
	maxCol, err := result.GetNumericColumn("price_max")
	if err != nil {
		t.Fatalf("expected price_max to be numeric: %v", err)
	}

	for i := 0; i < result.Len(); i++ {
		cat := catCol.AtAny(i).(string)

		if got := sumCol.AtAny(i).(float64); got != wantSum[cat] {
			t.Errorf("category %s: expected sum %v, got %v", cat, wantSum[cat], got)
		}
		if got := meanCol.AtAny(i).(float64); got != wantMean[cat] {
			t.Errorf("category %s: expected mean %v, got %v", cat, wantMean[cat], got)
		}
		if got := minCol.AtAny(i).(float64); got != wantMin[cat] {
			t.Errorf("category %s: expected min %v, got %v", cat, wantMin[cat], got)
		}
		if got := maxCol.AtAny(i).(float64); got != wantMax[cat] {
			t.Errorf("category %s: expected max %v, got %v", cat, wantMax[cat], got)
		}
	}

	t.Run("stddev of group B matches manual calculation", func(t *testing.T) {
		stddevCol, err := result.GetNumericColumn("price_stddev")
		if err != nil {
			t.Fatalf("expected price_stddev to be numeric: %v", err)
		}

		// group B = [30, 40, 50], mean 40, sample stddev (dof=1) = 10
		for i := 0; i < result.Len(); i++ {
			if catCol.AtAny(i).(string) == "B" {
				got := stddevCol.AtAny(i).(float64)
				if math.Abs(got-10) > 1e-9 {
					t.Errorf("expected group B stddev 10, got %v", got)
				}
			}
		}
	})

	t.Run("errors with no aggregations specified", func(t *testing.T) {
		if _, err := gb.AggregateColumns(map[string]map[string]dataframe.AggFunc{}); err == nil {
			t.Fatal("expected error for empty aggregations map")
		}
	})

	t.Run("errors for an unknown source column", func(t *testing.T) {
		_, err := gb.AggregateColumns(map[string]map[string]dataframe.AggFunc{
			"missing": {"sum": dataframe.AggSum},
		})
		if err == nil {
			t.Fatal("expected error aggregating a missing column")
		}
	})
}

func TestAggregateColumnsNonNumeric(t *testing.T) {
	df := newGroupedDataFrame(t)
	gb, err := df.GroupBy("category")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := gb.AggregateColumns(map[string]map[string]dataframe.AggFunc{
		"product": {"first": dataframe.AggFirst, "last": dataframe.AggLast, "count": dataframe.AggCount},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	firstCol, err := result.GetColumn("product_first")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lastCol, err := result.GetColumn("product_last")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	countCol, err := result.GetColumn("product_count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	catCol, _ := result.GetColumn("category")
	for i := 0; i < result.Len(); i++ {
		switch catCol.AtAny(i).(string) {
		case "A":
			if firstCol.AtAny(i) != "Widget" || lastCol.AtAny(i) != "Gadget" {
				t.Errorf("group A: expected first=Widget last=Gadget, got first=%v last=%v", firstCol.AtAny(i), lastCol.AtAny(i))
			}
			if countCol.AtAny(i) != 2 {
				t.Errorf("group A: expected count 2, got %v", countCol.AtAny(i))
			}
		case "B":
			if firstCol.AtAny(i) != "Foo" || lastCol.AtAny(i) != "Baz" {
				t.Errorf("group B: expected first=Foo last=Baz, got first=%v last=%v", firstCol.AtAny(i), lastCol.AtAny(i))
			}
			if countCol.AtAny(i) != 3 {
				t.Errorf("group B: expected count 3, got %v", countCol.AtAny(i))
			}
		}
	}
}

// Regression: Count() discards the error from gb.df.GetColumn(gb.groupColumn),
// unlike every other method in this file. That column lookup can genuinely
// fail: RenameColumn mutates the DataFrame in place (unlike every other
// DataFrame method, which returns a new DataFrame), so a DataFrameGroupBy's
// cached groupColumn name goes stale if the group column is renamed after
// grouping - and Count() then panics on a nil interface instead of
// returning a clear error.
func TestGroupByCount_GroupColumnRenamedAfterGrouping(t *testing.T) {
	df := newGroupedDataFrame(t)
	gb, err := df.GroupBy("category")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := df.RenameColumn("category", "cat"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := gb.Count(); err == nil {
		t.Fatal("expected Count() to return an error when its group column no longer exists, not panic or silently succeed")
	}
}

func TestGroupByCount(t *testing.T) {
	df := newGroupedDataFrame(t)
	gb, err := df.GroupBy("category")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := gb.Count()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rows, cols := result.Shape(); rows != 2 || cols != 2 {
		t.Fatalf("expected shape (2, 2), got (%d, %d)", rows, cols)
	}

	catCol, err := result.GetColumn("category")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	countCol, err := result.GetColumn("count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]int{"A": 2, "B": 3}
	for i := 0; i < result.Len(); i++ {
		cat := catCol.AtAny(i).(string)
		if countCol.AtAny(i) != want[cat] {
			t.Errorf("category %s: expected count %d, got %v", cat, want[cat], countCol.AtAny(i))
		}
	}
}
