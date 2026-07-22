package dataframe_test

import (
	"testing"

	dataframe "pango/Dataframe"
	"pango/series"
)

func newSampleDataFrame(t *testing.T) *dataframe.DataFrame {
	t.Helper()

	names := series.NewIndexSeries("name", []string{"Alice", "Bob", "Charlie", "David"})
	sales := series.NewIndexNumericSeries("sales", []int{100, 150, 200, 175})
	revenue := series.NewIndexNumericSeries("revenue", []float64{1500.5, 2250.75, 3000.0, 2625.25})

	df, err := dataframe.NewDataFrame(names, sales, revenue)
	if err != nil {
		t.Fatalf("unexpected error building sample DataFrame: %v", err)
	}
	return df
}

func TestNewDataFrame(t *testing.T) {
	t.Run("builds a DataFrame from series of equal length", func(t *testing.T) {
		df := newSampleDataFrame(t)

		rows, cols := df.Shape()
		if rows != 4 || cols != 3 {
			t.Errorf("expected shape (4, 3), got (%d, %d)", rows, cols)
		}
		if df.Len() != 4 {
			t.Errorf("expected Len() 4, got %d", df.Len())
		}

		wantCols := []string{"name", "sales", "revenue"}
		gotCols := df.Columns()
		if len(gotCols) != len(wantCols) {
			t.Fatalf("expected %d columns, got %d", len(wantCols), len(gotCols))
		}
		for i, c := range wantCols {
			if gotCols[i] != c {
				t.Errorf("expected column %d to be %q, got %q", i, c, gotCols[i])
			}
		}
	})

	t.Run("errors with no series", func(t *testing.T) {
		if _, err := dataframe.NewDataFrame(); err == nil {
			t.Fatal("expected error when constructing DataFrame with no series")
		}
	})

	t.Run("errors on mismatched lengths", func(t *testing.T) {
		a := series.NewIndexNumericSeries("a", []int{1, 2, 3})
		b := series.NewIndexNumericSeries("b", []int{1, 2})

		if _, err := dataframe.NewDataFrame(a, b); err == nil {
			t.Fatal("expected error for mismatched series lengths")
		}
	})

	t.Run("errors on duplicate column names", func(t *testing.T) {
		a := series.NewIndexNumericSeries("a", []int{1, 2, 3})
		b := series.NewIndexNumericSeries("a", []int{4, 5, 6})

		if _, err := dataframe.NewDataFrame(a, b); err == nil {
			t.Fatal("expected error for duplicate column names")
		}
	})
}

func TestNewDataFrameWithIndex(t *testing.T) {
	labels := []any{"x", "y", "z"}
	values := series.NewIndexNumericSeries("v", []int{1, 2, 3})

	df, err := dataframe.NewDataFrameWithIndex(labels, values)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotIndex := df.Index()
	for i, want := range labels {
		if gotIndex[i] != want {
			t.Errorf("expected index[%d] = %v, got %v", i, want, gotIndex[i])
		}
	}

	t.Run("errors when index length does not match series length", func(t *testing.T) {
		if _, err := dataframe.NewDataFrameWithIndex([]any{"x", "y"}, values); err == nil {
			t.Fatal("expected error for mismatched index/series length")
		}
	})
}

func TestGetColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	col, err := df.GetColumn("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col.Name() != "name" {
		t.Errorf("expected column name %q, got %q", "name", col.Name())
	}

	if _, err := df.GetColumn("missing"); err == nil {
		t.Fatal("expected error for missing column")
	}
}

func TestGetNumericColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	numCol, err := df.GetNumericColumn("sales")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if numCol.SumFloat() != 625 {
		t.Errorf("expected sales sum 625, got %v", numCol.SumFloat())
	}

	t.Run("errors for a non-numeric column", func(t *testing.T) {
		if _, err := df.GetNumericColumn("name"); err == nil {
			t.Fatal("expected error requesting a non-numeric column as numeric")
		}
	})

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.GetNumericColumn("missing"); err == nil {
			t.Fatal("expected error for missing column")
		}
	})
}

func TestAddColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	bonus := series.NewIndexNumericSeries("bonus", []float64{10, 20, 30, 40})
	if err := df.AddColumn(bonus); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, cols := df.Shape(); cols != 4 {
		t.Errorf("expected 4 columns after AddColumn, got %d", cols)
	}

	t.Run("errors on length mismatch", func(t *testing.T) {
		bad := series.NewIndexNumericSeries("bad", []float64{1, 2})
		if err := df.AddColumn(bad); err == nil {
			t.Fatal("expected error for length mismatch")
		}
	})

	t.Run("errors on duplicate column name", func(t *testing.T) {
		dup := series.NewIndexNumericSeries("bonus", []float64{1, 2, 3, 4})
		if err := df.AddColumn(dup); err == nil {
			t.Fatal("expected error for duplicate column name")
		}
	})
}

func TestDropColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	dropped, err := df.DropColumn("revenue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, cols := dropped.Shape(); cols != 2 {
		t.Errorf("expected 2 columns after drop, got %d", cols)
	}
	if _, cols := df.Shape(); cols != 3 {
		t.Errorf("original DataFrame should be unmodified, expected 3 columns, got %d", cols)
	}
	if _, err := dropped.GetColumn("revenue"); err == nil {
		t.Fatal("expected dropped column to be gone")
	}

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.DropColumn("missing"); err == nil {
			t.Fatal("expected error for missing column")
		}
	})
}

func TestRenameColumn(t *testing.T) {
	sales := series.NewIndexNumericSeries("sales", []int{1, 2, 3})
	df, err := dataframe.NewDataFrame(sales)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := df.RenameColumn("sales", "units"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := df.GetColumn("units"); err != nil {
		t.Fatalf("expected renamed column %q to exist: %v", "units", err)
	}
	if _, err := df.GetColumn("sales"); err == nil {
		t.Fatal("expected old column name to be gone after rename")
	}

	// Regression: RenameColumn must not mutate the caller's original series,
	// since DataFrame stores columns by reference rather than by copy.
	if sales.Name() != "sales" {
		t.Errorf("RenameColumn mutated the caller's original series; name is now %q", sales.Name())
	}

	t.Run("errors renaming a missing column", func(t *testing.T) {
		if err := df.RenameColumn("missing", "x"); err == nil {
			t.Fatal("expected error renaming a missing column")
		}
	})

	t.Run("errors when the new name already exists", func(t *testing.T) {
		other := series.NewIndexNumericSeries("other", []int{1, 2, 3})
		df2, _ := dataframe.NewDataFrame(series.NewIndexNumericSeries("units", []int{1, 2, 3}), other)
		if err := df2.RenameColumn("other", "units"); err == nil {
			t.Fatal("expected error when renaming to an already-existing column name")
		}
	})
}

func TestSelect(t *testing.T) {
	df := newSampleDataFrame(t)

	selected, err := df.Select("name", "sales")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, cols := selected.Shape(); cols != 2 {
		t.Errorf("expected 2 columns, got %d", cols)
	}
	if _, err := selected.GetColumn("revenue"); err == nil {
		t.Fatal("expected unselected column to be absent")
	}

	t.Run("errors with no columns specified", func(t *testing.T) {
		if _, err := df.Select(); err == nil {
			t.Fatal("expected error selecting zero columns")
		}
	})

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.Select("missing"); err == nil {
			t.Fatal("expected error selecting a missing column")
		}
	})
}

func TestHeadTail(t *testing.T) {
	df := newSampleDataFrame(t)

	head := df.Head(2)
	if head.Len() != 2 {
		t.Errorf("expected Head(2) to have 2 rows, got %d", head.Len())
	}

	tail := df.Tail(2)
	if tail.Len() != 2 {
		t.Errorf("expected Tail(2) to have 2 rows, got %d", tail.Len())
	}

	t.Run("defaults to 5 rows when n <= 0", func(t *testing.T) {
		got := df.Head(0)
		if got.Len() != df.Len() {
			t.Errorf("expected Head(0) to default to 5 (clamped to %d rows), got %d", df.Len(), got.Len())
		}
	})

	t.Run("clamps n to available rows", func(t *testing.T) {
		got := df.Head(100)
		if got.Len() != df.Len() {
			t.Errorf("expected Head(100) to clamp to %d rows, got %d", df.Len(), got.Len())
		}
	})

	// Regression: Head/Tail used to always downgrade columns to a
	// non-numeric fallback series, breaking GetNumericColumn on the result.
	t.Run("preserves numeric column capability", func(t *testing.T) {
		if _, err := head.GetNumericColumn("sales"); err != nil {
			t.Errorf("expected Head result to preserve numeric column: %v", err)
		}
		if _, err := tail.GetNumericColumn("revenue"); err != nil {
			t.Errorf("expected Tail result to preserve numeric column: %v", err)
		}
	})

	t.Run("tail returns the last rows in original order", func(t *testing.T) {
		nameCol, err := tail.GetColumn("name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nameCol.AtAny(0) != "Charlie" || nameCol.AtAny(1) != "David" {
			t.Errorf("expected tail rows [Charlie, David], got [%v, %v]", nameCol.AtAny(0), nameCol.AtAny(1))
		}
	})
}

func TestFilterColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	filtered, err := df.FilterColumn("sales", func(v any) bool {
		return v.(int) > 150
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filtered.Len() != 2 {
		t.Errorf("expected 2 matching rows, got %d", filtered.Len())
	}

	// Regression: a predicate matching zero rows used to return an error
	// instead of a valid empty DataFrame.
	t.Run("returns an empty DataFrame, not an error, when nothing matches", func(t *testing.T) {
		empty, err := df.FilterColumn("sales", func(v any) bool {
			return v.(int) > 100000
		})
		if err != nil {
			t.Fatalf("expected no error for zero matches, got: %v", err)
		}
		if empty.Len() != 0 {
			t.Errorf("expected 0 rows, got %d", empty.Len())
		}
	})

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.FilterColumn("missing", func(v any) bool { return true }); err == nil {
			t.Fatal("expected error for missing column")
		}
	})
}

func TestResetIndex(t *testing.T) {
	filtered, err := newSampleDataFrame(t).FilterColumn("sales", func(v any) bool { return v.(int) > 150 })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reset := filtered.ResetIndex()
	gotIndex := reset.Index()
	for i, idx := range gotIndex {
		if idx != i {
			t.Errorf("expected reset index[%d] = %d, got %v", i, i, idx)
		}
	}
}

func TestSetIndexFromColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	indexed, err := df.SetIndexFromColumn("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := indexed.GetColumn("name"); err == nil {
		t.Fatal("expected the index source column to be removed")
	}

	gotIndex := indexed.Index()
	if gotIndex[0] != "Alice" {
		t.Errorf("expected index[0] = Alice, got %v", gotIndex[0])
	}

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.SetIndexFromColumn("missing"); err == nil {
			t.Fatal("expected error for missing column")
		}
	})
}

func TestApplyToColumn(t *testing.T) {
	df := newSampleDataFrame(t)

	doubled, err := df.ApplyToColumn("sales", func(v any) any {
		return v.(int) * 2
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	col, err := doubled.GetNumericColumn("sales")
	if err != nil {
		t.Fatalf("expected transformed column to stay numeric: %v", err)
	}
	if col.SumFloat() != 1250 {
		t.Errorf("expected doubled sum 1250, got %v", col.SumFloat())
	}

	// Untouched columns must be unaffected.
	if _, err := doubled.GetNumericColumn("revenue"); err != nil {
		t.Errorf("expected untouched revenue column to remain numeric: %v", err)
	}

	t.Run("errors for a missing column", func(t *testing.T) {
		if _, err := df.ApplyToColumn("missing", func(v any) any { return v }); err == nil {
			t.Fatal("expected error for missing column")
		}
	})
}

func TestApplyToColumns(t *testing.T) {
	df := newSampleDataFrame(t)

	withTotal, err := df.ApplyToColumns([]string{"sales", "revenue"}, func(row map[string]any) any {
		return float64(row["sales"].(int)) + row["revenue"].(float64)
	}, "combined")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	combined, err := withTotal.GetNumericColumn("combined")
	if err != nil {
		t.Fatalf("expected result column to be numeric: %v", err)
	}
	if combined.AtAny(0) != 1600.5 {
		t.Errorf("expected combined[0] = 1600.5, got %v", combined.AtAny(0))
	}

	t.Run("errors for a missing source column", func(t *testing.T) {
		_, err := df.ApplyToColumns([]string{"missing"}, func(row map[string]any) any { return nil }, "x")
		if err == nil {
			t.Fatal("expected error for missing source column")
		}
	})
}

func TestRange(t *testing.T) {
	df := newSampleDataFrame(t)

	var names []string
	for _, row := range df.Range() {
		names = append(names, row["name"].(string))
	}

	want := []string{"Alice", "Bob", "Charlie", "David"}
	if len(names) != len(want) {
		t.Fatalf("expected %d rows, got %d", len(want), len(names))
	}
	for i, n := range want {
		if names[i] != n {
			t.Errorf("expected row %d name %q, got %q", i, n, names[i])
		}
	}
}

// Regression: DataFrame construction stores the caller's Series by
// reference instead of copying it. If the caller keeps the typed pointer
// and mutates it afterwards (e.g. via Append), the DataFrame's column
// silently desyncs from DataFrame.Len()/index without any error.
func TestNewDataFrame_DoesNotDesyncAfterCallerMutatesOriginalSeries(t *testing.T) {
	nums := series.NewIndexNumericSeries("a", []int{1, 2, 3})
	df, err := dataframe.NewDataFrame(nums)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Caller appends more data to the series it already handed to the
	// DataFrame - this should not be visible through df's column at all,
	// since the DataFrame is supposed to own an independent copy.
	more := series.NewSeries("a", []int{4, 5}, []int{3, 4})
	nums.Append(more)

	col, err := df.GetColumn("a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if col.Len() != df.Len() {
		t.Fatalf("DataFrame column length (%d) desynced from DataFrame.Len() (%d) after external mutation of the original series", col.Len(), df.Len())
	}
	if col.Len() != 3 {
		t.Fatalf("expected DataFrame's column to stay at 3 elements (unaffected by the caller's later Append), got %d", col.Len())
	}
}

func TestColumnTypesAndStringers(t *testing.T) {
	df := newSampleDataFrame(t)

	types := df.ColumnTypes()
	if len(types) != 3 {
		t.Errorf("expected 3 column types, got %d", len(types))
	}

	// String() and Info() are human-readable helpers; just make sure they
	// don't panic and mention the DataFrame's shape.
	if s := df.String(); s == "" {
		t.Error("expected non-empty String() output")
	}
	if s := df.Info(); s == "" {
		t.Error("expected non-empty Info() output")
	}
}
