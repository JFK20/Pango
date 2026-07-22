package series_test

import (
	"testing"

	series "pango/series"
)

func TestNewSeries(t *testing.T) {
	t.Run("creates series with valid data", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		index := []string{"a", "b", "c", "d", "e"}
		s := series.NewSeries("test", values, index)

		if s == nil {
			t.Fatal("expected series to be created")
		}
		if s.Name() != "test" {
			t.Errorf("expected name 'test', got %s", s.Name())
		}
		if s.Len() != 5 {
			t.Errorf("expected 5 values, got %d", s.Len())
		}
	})

	t.Run("panics with empty values", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty values")
			}
		}()
		series.NewSeries[int, string]("test", []int{}, nil)
	})
}

func TestNewSeries_EdgeCases(t *testing.T) {
	t.Run("panics with nil index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for nil index")
			}
		}()
		values := []int{1, 2, 3}
		series.NewSeries[int, string]("test", values, nil)
	})

	t.Run("panics with mismatched lengths", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for mismatched lengths")
			}
		}()
		values := []int{1, 2, 3}
		index := []string{"a", "b"} // wrong length
		series.NewSeries("test", values, index)
	})
}

func TestLen(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	s := series.NewIndexSeries("test", values)

	if s.Len() != 5 {
		t.Errorf("expected length 5, got %d", s.Len())
	}
}

func TestName(t *testing.T) {
	values := []string{"a", "b", "c"}
	s := series.NewIndexSeries("my_series", values)

	if s.Name() != "my_series" {
		t.Errorf("expected name 'my_series', got %s", s.Name())
	}
}

func TestSetName(t *testing.T) {
	values := []int{1, 2, 3}
	s := series.NewIndexSeries("old_name", values)
	s.SetName("new_name")

	if s.Name() != "new_name" {
		t.Errorf("expected name 'new_name', got %s", s.Name())
	}
}

func TestValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	s := series.NewIndexSeries("test", values)

	copied := s.Values()

	// Check values are correct
	for i, v := range copied {
		if v != values[i] {
			t.Errorf("expected value %d at index %d, got %d", values[i], i, v)
		}
	}

	// Verify it's a copy by modifying it
	copied[0] = 999
	if s.At(0) == 999 {
		t.Error("Values() should return a copy, not the original slice")
	}
}

func TestIndex(t *testing.T) {
	t.Run("returns provided index", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		idx := s.Index()
		if len(idx) != 3 {
			t.Errorf("expected index length 3, got %d", len(idx))
		}
	})

	t.Run("generates index when nil", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("test", values)

		idx := s.Index()
		if len(idx) != 3 {
			t.Errorf("expected index length 3, got %d", len(idx))
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("gets value at valid index", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		label, value := s.AtIndex(1)
		if label != "b" {
			t.Errorf("expected label 'b', got %s", label)
		}
		if value != 20 {
			t.Errorf("expected value 20, got %v", value)
		}
	})

	t.Run("panics on negative index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative index")
			}
		}()
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("test", values)
		s.At(-1)
	})

	t.Run("panics on out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds index")
			}
		}()
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("test", values)
		s.At(10)
	})

	t.Run("panics on out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds index")
			}
		}()
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("test", values)
		s.AtIndex(3)
	})
}

func TestSeries_IndexGet(t *testing.T) {
	values := []int{1, 2, 3}
	index := []string{"a", "b", "c"}
	s := series.NewSeries("test", values, index)

	val := s.Get("b")
	if val != 2 {
		t.Errorf("expected value 2 for index 'b', got %v", val)
	}
}

func TestString(t *testing.T) {
	t.Run("prints series with name", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test_series", values, index)

		str := s.String()
		if str == "" {
			t.Error("expected non-empty string")
		}
		// Should contain the name
		if len(str) < len("test_series") {
			t.Error("expected string to contain series name")
		}
	})

	t.Run("handles series with more than 10 elements", func(t *testing.T) {
		values := make([]int, 15)
		for i := range values {
			values[i] = i
		}
		s := series.NewIndexSeries("long_series", values)

		str := s.String()
		if str == "" {
			t.Error("expected non-empty string")
		}
	})
}

func TestHead(t *testing.T) {
	t.Run("returns first n elements", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		index := []string{"a", "b", "c", "d", "e"}
		s := series.NewSeries("test", values, index)

		head := s.Head(3)
		if head.Len() != 3 {
			t.Errorf("expected length 3, got %d", head.Len())
		}
		if head.At(0) != 1 || head.At(1) != 2 || head.At(2) != 3 {
			t.Error("head values are incorrect")
		}
	})

	t.Run("returns all elements when n exceeds length", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("test", values)

		head := s.Head(10)
		if head.Len() != 3 {
			t.Errorf("expected length 3, got %d", head.Len())
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		s := series.NewIndexSeries("my_series", values)

		head := s.Head(2)
		if head.Name() != "my_series" {
			t.Errorf("expected name 'my_series', got %s", head.Name())
		}
	})
}

func TestTail(t *testing.T) {
	t.Run("returns last n elements", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		index := []string{"a", "b", "c", "d", "e"}
		s := series.NewSeries("test", values, index)

		tail := s.Tail(3)
		if tail.Len() != 3 {
			t.Errorf("expected length 3, got %d", tail.Len())
		}
		if tail.At(0) != 3 || tail.At(1) != 4 || tail.At(2) != 5 {
			t.Error("tail values are incorrect")
		}
	})

	t.Run("returns all elements when n exceeds length", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("test", values)

		tail := s.Tail(10)
		if tail.Len() != 3 {
			t.Errorf("expected length 3, got %d", tail.Len())
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		s := series.NewIndexSeries("my_series", values)

		tail := s.Tail(2)
		if tail.Name() != "my_series" {
			t.Errorf("expected name 'my_series', got %s", tail.Name())
		}
	})
}

func TestAppend(t *testing.T) {
	t.Run("appends Series", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := series.NewIndexSeries("series1", values)

		values = []int{4, 5, 6}
		o := series.NewIndexSeries("series2", values)

		s.Append(o)
		if s.Len() != 6 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.At(3) != 4 {
			t.Errorf("expected value to be 4, got %d", s.At(3))
		}
	})
}

func TestPrepend(t *testing.T) {
	t.Run("prepends Series to the beginning", func(t *testing.T) {
		values := []int{4, 5, 6}
		s := series.NewIndexSeries("series1", values)

		values = []int{1, 2, 3}
		o := series.NewIndexSeries("series2", values)

		s.Prepend(o)
		if s.Len() != 6 {
			t.Errorf("expected length 6, got %d", s.Len())
		}
		if s.At(0) != 1 {
			t.Errorf("expected first value to be 1, got %d", s.At(0))
		}
		if s.At(2) != 3 {
			t.Errorf("expected third value to be 3, got %d", s.At(2))
		}
		if s.At(3) != 4 {
			t.Errorf("expected fourth value to be 4, got %d", s.At(3))
		}
	})
}

func TestGet_NonExisting(t *testing.T) {
	t.Run("panics for non-existing label", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for not existing index")
			}
		}()
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		s.Get("nonexistent")
	})
}

func TestResetIndex(t *testing.T) {
	t.Run("resets index to 0..n-1", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		resetSeries := s.ResetIndex()
		if resetSeries.Len() != 3 {
			t.Errorf("expected length 3, got %d", resetSeries.Len())
		}

		idx := resetSeries.Index()
		for i := range idx {
			if idx[i] != i {
				t.Errorf("expected index %d at position %d, got %d", i, i, idx[i])
			}
		}

		// Check that values are preserved
		if resetSeries.At(0) != 10 || resetSeries.At(1) != 20 || resetSeries.At(2) != 30 {
			t.Error("values should be preserved after ResetIndex")
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"x", "y", "z"}
		s := series.NewSeries("my_series", values, index)

		resetSeries := s.ResetIndex()
		if resetSeries.Name() != "my_series" {
			t.Errorf("expected name 'my_series', got %s", resetSeries.Name())
		}
	})
}

func TestSetIndex(t *testing.T) {
	t.Run("sets new index", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		newIndex := []float64{1.1, 2.2, 3.3}
		newSeries := series.SetIndex(s, newIndex)

		if newSeries.Len() != 3 {
			t.Errorf("expected length 3, got %d", newSeries.Len())
		}

		idx := newSeries.Index()
		if idx[0] != 1.1 || idx[1] != 2.2 || idx[2] != 3.3 {
			t.Error("new index not set correctly")
		}

		// Check that values are preserved
		if newSeries.At(0) != 10 || newSeries.At(1) != 20 || newSeries.At(2) != 30 {
			t.Error("values should be preserved after SetIndex")
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("my_series", values, index)

		newIndex := []int{10, 20, 30}
		newSeries := series.SetIndex(s, newIndex)
		if newSeries.Name() != "my_series" {
			t.Errorf("expected name 'my_series', got %s", newSeries.Name())
		}
	})

	t.Run("panics with mismatched index length", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for mismatched index length")
			}
		}()
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		newIndex := []int{1, 2} // wrong length
		series.SetIndex(s, newIndex)
	})
}

func TestIndex_Copy(t *testing.T) {
	t.Run("Index returns a copy not the original", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		copiedIndex := s.Index()
		copiedIndex[0] = "modified"

		// Verify original index is unchanged
		originalIndex := s.Index()
		if originalIndex[0] != "a" {
			t.Error("Index() should return a copy, not the original slice")
		}
	})
}

func TestSortByIndex(t *testing.T) {
	t.Run("sorts series in ascending order", func(t *testing.T) {
		values := []int{30, 10, 20}
		index := []string{"c", "a", "b"}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByIndex(s, true)
		sortedIndex := sorted.Index()
		if sortedIndex[0] != "a" || sortedIndex[1] != "b" || sortedIndex[2] != "c" {
			t.Error("series not sorted correctly in ascending order")
		}
	})

	t.Run("sorts series in descending order", func(t *testing.T) {
		values := []int{30, 10, 20}
		index := []string{"c", "a", "b"}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByIndex(s, false)
		sortedIndex := sorted.Index()
		if sortedIndex[0] != "c" || sortedIndex[1] != "b" || sortedIndex[2] != "a" {
			t.Error("series not sorted correctly in descending order")
		}
	})
}

func TestSeries_IsIn(t *testing.T) {
	t.Run("returns true if series contains index", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		if !s.IsIn(20) {
			t.Error("expected IsIn to return true for existing value")
		}
	})

	t.Run("returns false if series does not contain index", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		if s.IsIn(40) {
			t.Error("expected IsIn to return false for non-existing value")
		}
	})
}

func TestSortByValue(t *testing.T) {
	t.Run("sorts by integer values ascending", func(t *testing.T) {
		values := []int{30, 10, 20}
		index := []string{"c", "a", "b"}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByValue(s, true)

		// Check values are sorted
		expectedValues := []int{10, 20, 30}
		expectedIndex := []string{"a", "b", "c"}

		vals := sorted.Values()
		for i, exp := range expectedValues {
			if vals[i] != exp {
				t.Errorf("expected value %d at position %d, got %d", exp, i, vals[i])
			}
		}

		// Check index is reordered accordingly
		idx := sorted.Index()
		for i, exp := range expectedIndex {
			if idx[i] != exp {
				t.Errorf("expected label %s at position %d, got %s", exp, i, idx[i])
			}
		}
	})

	t.Run("sorts by integer values descending", func(t *testing.T) {
		values := []int{30, 10, 20}
		index := []string{"c", "a", "b"}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByValue(s, false)

		// Check values are sorted descending
		expectedValues := []int{30, 20, 10}
		expectedIndex := []string{"c", "b", "a"}

		vals := sorted.Values()
		for i, exp := range expectedValues {
			if vals[i] != exp {
				t.Errorf("expected value %d at position %d, got %d", exp, i, vals[i])
			}
		}

		// Check index is reordered accordingly
		idx := sorted.Index()
		for i, exp := range expectedIndex {
			if idx[i] != exp {
				t.Errorf("expected label %s at position %d, got %s", exp, i, idx[i])
			}
		}
	})

	t.Run("sorts by string values ascending", func(t *testing.T) {
		values := []string{"zebra", "apple", "mango"}
		index := []int{3, 1, 2}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByValue(s, true)

		// Check values are sorted
		expectedValues := []string{"apple", "mango", "zebra"}
		expectedIndex := []int{1, 2, 3}

		vals := sorted.Values()
		for i, exp := range expectedValues {
			if vals[i] != exp {
				t.Errorf("expected value %s at position %d, got %s", exp, i, vals[i])
			}
		}

		// Check index is reordered accordingly
		idx := sorted.Index()
		for i, exp := range expectedIndex {
			if idx[i] != exp {
				t.Errorf("expected label %d at position %d, got %d", exp, i, idx[i])
			}
		}
	})

	t.Run("sorts by string values descending", func(t *testing.T) {
		values := []string{"zebra", "apple", "mango"}
		index := []int{3, 1, 2}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByValue(s, false)

		// Check values are sorted descending
		expectedValues := []string{"zebra", "mango", "apple"}
		expectedIndex := []int{3, 2, 1}

		vals := sorted.Values()
		for i, exp := range expectedValues {
			if vals[i] != exp {
				t.Errorf("expected value %s at position %d, got %s", exp, i, vals[i])
			}
		}

		// Check index is reordered accordingly
		idx := sorted.Index()
		for i, exp := range expectedIndex {
			if idx[i] != exp {
				t.Errorf("expected label %d at position %d, got %d", exp, i, idx[i])
			}
		}
	})

	t.Run("original series is not modified", func(t *testing.T) {
		values := []int{30, 10, 20}
		index := []string{"c", "a", "b"}
		s := series.NewSeries("test", values, index)

		originalFirst := s.At(0)
		series.SortByValue(s, true)

		if s.At(0) != originalFirst {
			t.Error("original series should not be modified")
		}
	})

	t.Run("handles series with duplicate values", func(t *testing.T) {
		values := []int{20, 10, 20, 10}
		index := []string{"a", "b", "c", "d"}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByValue(s, true)

		// Check that all values are present
		if sorted.Len() != 4 {
			t.Errorf("expected length 4, got %d", sorted.Len())
		}

		// Check values are sorted
		vals := sorted.Values()
		if vals[0] != 10 || vals[1] != 10 || vals[2] != 20 || vals[3] != 20 {
			t.Error("values not sorted correctly with duplicates")
		}
	})

	t.Run("handles single element series", func(t *testing.T) {
		values := []int{42}
		index := []string{"answer"}
		s := series.NewSeries("test", values, index)

		sorted := series.SortByValue(s, true)

		if sorted.Len() != 1 {
			t.Errorf("expected length 1, got %d", sorted.Len())
		}
		if sorted.At(0) != 42 {
			t.Errorf("expected value 42, got %d", sorted.At(0))
		}
	})
}

func TestCopy(t *testing.T) {
	t.Run("creates a deep copy of the series", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		index := []string{"a", "b", "c", "d", "e"}
		s := series.NewSeries("test", values, index)

		copied := s.Copy()

		// Check that the copy has the same length
		if copied.Len() != s.Len() {
			t.Errorf("expected copied series length %d, got %d", s.Len(), copied.Len())
		}

		// Check that the copy has the same name
		if copied.Name() != s.Name() {
			t.Errorf("expected copied series name '%s', got '%s'", s.Name(), copied.Name())
		}

		// Check that values are the same
		for i := 0; i < s.Len(); i++ {
			if copied.At(i) != s.At(i) {
				t.Errorf("expected value %d at position %d, got %d", s.At(i), i, copied.At(i))
			}
		}

		// Check that index is the same
		originalIndex := s.Index()
		copiedIndex := copied.Index()
		for i := 0; i < s.Len(); i++ {
			if copiedIndex[i] != originalIndex[i] {
				t.Errorf("expected index %s at position %d, got %s", originalIndex[i], i, copiedIndex[i])
			}
		}

		// check that modifying the copy does not affect the original
		copied.SetName("modified")
		if s.Name() == "modified" {
			t.Error("modifying copy name should not affect original series name")
		}

		// check that modifying the copy values does not affect the original
		copiedValues := copied.Values()
		copiedValues[0] = 999
		if s.At(0) == 999 {
			t.Error("modifying copy values should not affect original series values")
		}

		// check that modifying the copy index does not affect the original
		copiedIndex = copied.Index()
		copiedIndex[0] = "modified"
		originalIndex = s.Index()
		if originalIndex[0] == "modified" {
			t.Error("modifying copy index should not affect original series index")
		}

		// check that modifying the original does not affect the copy
		s.SetName("original_modified")
		if copied.Name() == "original_modified" {
			t.Error("modifying original name should not affect copy series name")
		}

		// Append mutates the original's backing values/index slices (see
		// Series.Append); this proves Copy() holds independent backing
		// arrays rather than just an independent struct wrapping the same
		// slices.
		copiedLenBeforeAppend := copied.Len()
		more := series.NewSeries("more", []int{888}, []string{"z"})
		s.Append(more)
		if copied.Len() != copiedLenBeforeAppend {
			t.Error("appending to the original after Copy() should not affect the copy's length")
		}
		if copied.At(0) == 888 {
			t.Error("modifying original values should not affect copy series values")
		}
	})
}

func TestCopyAny(t *testing.T) {
	t.Run("creates a copy through interface", func(t *testing.T) {
		values := []string{"apple", "banana", "cherry"}
		index := []int{10, 20, 30}
		s := series.NewSeries("fruits", values, index)

		copied := s.CopyAny()

		if copied == nil {
			t.Fatal("expected copy to be created")
		}

		if copied.Name() != s.Name() {
			t.Errorf("expected name %s, got %s", s.Name(), copied.Name())
		}

		if copied.Len() != s.Len() {
			t.Errorf("expected length %d, got %d", s.Len(), copied.Len())
		}
	})
}

func TestValuesAny(t *testing.T) {
	t.Run("returns values as []any for int series", func(t *testing.T) {
		values := []int{1, 2, 3, 4}
		index := []string{"a", "b", "c", "d"}
		s := series.NewSeries("test", values, index)

		result := s.ValuesAny()

		if len(result) != len(values) {
			t.Errorf("expected length %d, got %d", len(values), len(result))
		}

		for i, v := range values {
			if result[i] != v {
				t.Errorf("expected value %v at index %d, got %v", v, i, result[i])
			}
		}
	})
}

func TestIndexAny(t *testing.T) {
	t.Run("returns index as []any for string index", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		result := s.IndexAny()

		if len(result) != len(index) {
			t.Errorf("expected length %d, got %d", len(index), len(result))
		}

		for i, idx := range index {
			if result[i] != idx {
				t.Errorf("expected index %v at position %d, got %v", idx, i, result[i])
			}
		}
	})
}

func TestAtAny(t *testing.T) {
	t.Run("returns value at index as any", func(t *testing.T) {
		values := []string{"first", "second", "third"}
		index := []int{1, 2, 3}
		s := series.NewSeries("test", values, index)

		result := s.AtAny(0)
		if result != "first" {
			t.Errorf("expected 'first', got %v", result)
		}

		result = s.AtAny(2)
		if result != "third" {
			t.Errorf("expected 'third', got %v", result)
		}
	})
}

func TestAtIndexAny(t *testing.T) {
	t.Run("returns label and value as any", func(t *testing.T) {
		values := []int{100, 200, 300}
		index := []string{"x", "y", "z"}
		s := series.NewSeries("test", values, index)

		label, value := s.AtIndexAny(0)
		if label != "x" {
			t.Errorf("expected label 'x', got %v", label)
		}
		if value != 100 {
			t.Errorf("expected value 100, got %v", value)
		}

		label, value = s.AtIndexAny(2)
		if label != "z" {
			t.Errorf("expected label 'z', got %v", label)
		}
		if value != 300 {
			t.Errorf("expected value 300, got %v", value)
		}
	})
}

func TestGetValueType(t *testing.T) {
	t.Run("returns correct type for int values", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		typeStr := s.GetValueType()
		if typeStr != "int" {
			t.Errorf("expected 'int', got %s", typeStr)
		}
	})
}

func TestGetIndexType(t *testing.T) {
	t.Run("returns correct type for string index", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := series.NewSeries("test", values, index)

		typeStr := s.GetIndexType()
		if typeStr != "string" {
			t.Errorf("expected 'string', got %s", typeStr)
		}
	})

	t.Run("returns correct type for int index", func(t *testing.T) {
		values := []string{"x", "y", "z"}
		index := []int{10, 20, 30}
		s := series.NewSeries("test", values, index)

		typeStr := s.GetIndexType()
		if typeStr != "int" {
			t.Errorf("expected 'int', got %s", typeStr)
		}
	})
}
