package series

import (
	"testing"
)

func TestNewSeries(t *testing.T) {
	t.Run("creates series with valid data", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		index := []string{"a", "b", "c", "d", "e"}
		s := NewSeries("test", values, index)

		if s == nil {
			t.Fatal("expected series to be created")
		}
		if s.name != "test" {
			t.Errorf("expected name 'test', got %s", s.name)
		}
		if len(s.values) != 5 {
			t.Errorf("expected 5 values, got %d", len(s.values))
		}
	})

	t.Run("panics with empty values", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty values")
			}
		}()
		NewSeries[int, string]("test", []int{}, nil)
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
		NewSeries[int, string]("test", values, nil)
	})

	t.Run("panics with mismatched lengths", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for mismatched lengths")
			}
		}()
		values := []int{1, 2, 3}
		index := []string{"a", "b"} // wrong length
		NewSeries("test", values, index)
	})
}

func TestLen(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	s := NewIndexSeries("test", values)

	if s.Len() != 5 {
		t.Errorf("expected length 5, got %d", s.Len())
	}
}

func TestName(t *testing.T) {
	values := []string{"a", "b", "c"}
	s := NewIndexSeries("my_series", values)

	if s.Name() != "my_series" {
		t.Errorf("expected name 'my_series', got %s", s.Name())
	}
}

func TestSetName(t *testing.T) {
	values := []int{1, 2, 3}
	s := NewIndexSeries("old_name", values)
	s.SetName("new_name")

	if s.Name() != "new_name" {
		t.Errorf("expected name 'new_name', got %s", s.Name())
	}
}

func TestValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	s := NewIndexSeries("test", values)

	copied := s.Values()

	// Check values are correct
	for i, v := range copied {
		if v != values[i] {
			t.Errorf("expected value %d at index %d, got %d", values[i], i, v)
		}
	}

	// Verify it's a copy by modifying it
	copied[0] = 999
	if s.values[0] == 999 {
		t.Error("Values() should return a copy, not the original slice")
	}
}

func TestIndex(t *testing.T) {
	t.Run("returns provided index", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := NewSeries("test", values, index)

		idx := s.Index()
		if len(idx) != 3 {
			t.Errorf("expected index length 3, got %d", len(idx))
		}
	})

	t.Run("generates index when nil", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := NewIndexSeries("test", values)

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
		s := NewSeries("test", values, index)

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
		s := NewIndexSeries("test", values)
		s.At(-1)
	})

	t.Run("panics on out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds index")
			}
		}()
		values := []int{1, 2, 3}
		s := NewIndexSeries("test", values)
		s.At(10)
	})

	t.Run("panics on out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds index")
			}
		}()
		values := []int{1, 2, 3}
		s := NewIndexSeries("test", values)
		s.AtIndex(3)
	})
}

func TestSeries_IndexGet(t *testing.T) {
	values := []int{1, 2, 3}
	index := []string{"a", "b", "c"}
	s := NewSeries("test", values, index)

	val := s.Get("b")
	if val != 2 {
		t.Errorf("expected value 2 for index 'b', got %v", val)
	}
}

func TestString(t *testing.T) {
	t.Run("prints series with name", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := NewSeries("test_series", values, index)

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
		s := NewIndexSeries("long_series", values)

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
		s := NewSeries("test", values, index)

		head := s.Head(3)
		if head.Len() != 3 {
			t.Errorf("expected length 3, got %d", head.Len())
		}
		if head.values[0] != 1 || head.values[1] != 2 || head.values[2] != 3 {
			t.Error("head values are incorrect")
		}
	})

	t.Run("returns all elements when n exceeds length", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := NewIndexSeries("test", values)

		head := s.Head(10)
		if head.Len() != 3 {
			t.Errorf("expected length 3, got %d", head.Len())
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		s := NewIndexSeries("my_series", values)

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
		s := NewSeries("test", values, index)

		tail := s.Tail(3)
		if tail.Len() != 3 {
			t.Errorf("expected length 3, got %d", tail.Len())
		}
		if tail.values[0] != 3 || tail.values[1] != 4 || tail.values[2] != 5 {
			t.Error("tail values are incorrect")
		}
	})

	t.Run("returns all elements when n exceeds length", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := NewIndexSeries("test", values)

		tail := s.Tail(10)
		if tail.Len() != 3 {
			t.Errorf("expected length 3, got %d", tail.Len())
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		s := NewIndexSeries("my_series", values)

		tail := s.Tail(2)
		if tail.Name() != "my_series" {
			t.Errorf("expected name 'my_series', got %s", tail.Name())
		}
	})
}

func TestAppend(t *testing.T) {
	t.Run("appends Series", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := NewIndexSeries("series1", values)

		values = []int{4, 5, 6}
		o := NewIndexSeries("series2", values)

		s.Append(o)
		if s.Len() != 6 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.values[3] != 4 {
			t.Errorf("expected value to be 4, got %d", s.values[3])
		}
	})
}

func TestPrepend(t *testing.T) {
	t.Run("prepends Series to the beginning", func(t *testing.T) {
		values := []int{4, 5, 6}
		s := NewIndexSeries("series1", values)

		values = []int{1, 2, 3}
		o := NewIndexSeries("series2", values)

		s.Prepend(o)
		if s.Len() != 6 {
			t.Errorf("expected length 6, got %d", s.Len())
		}
		if s.values[0] != 1 {
			t.Errorf("expected first value to be 1, got %d", s.values[0])
		}
		if s.values[2] != 3 {
			t.Errorf("expected third value to be 3, got %d", s.values[2])
		}
		if s.values[3] != 4 {
			t.Errorf("expected fourth value to be 4, got %d", s.values[3])
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
		s := NewSeries("test", values, index)

		s.Get("nonexistent")
	})
}

func TestResetIndex(t *testing.T) {
	t.Run("resets index to 0..n-1", func(t *testing.T) {
		values := []int{10, 20, 30}
		index := []string{"a", "b", "c"}
		s := NewSeries("test", values, index)

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
		s := NewSeries("my_series", values, index)

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
		s := NewSeries("test", values, index)

		newIndex := []float64{1.1, 2.2, 3.3}
		newSeries := SetIndex(s, newIndex)

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
		s := NewSeries("my_series", values, index)

		newIndex := []int{10, 20, 30}
		newSeries := SetIndex(s, newIndex)
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
		s := NewSeries("test", values, index)

		newIndex := []int{1, 2} // wrong length
		SetIndex(s, newIndex)
	})
}

func TestIndex_Copy(t *testing.T) {
	t.Run("Index returns a copy not the original", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		s := NewSeries("test", values, index)

		copiedIndex := s.Index()
		copiedIndex[0] = "modified"

		// Verify original index is unchanged
		originalIndex := s.Index()
		if originalIndex[0] != "a" {
			t.Error("Index() should return a copy, not the original slice")
		}
	})
}

// we cant create an empty series using NewSeries as it panics
func TestIsEmpty(t *testing.T) {
	t.Run("returns false for non-empty series", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := NewIndexSeries("test", values)

		if s.isEmpty() {
			t.Error("expected isEmpty to return false for non-empty series")
		}
	})
}
