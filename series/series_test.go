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

func TestLen(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	s := NewSeriesInt("test", values)

	if s.Len() != 5 {
		t.Errorf("expected length 5, got %d", s.Len())
	}
}

func TestName(t *testing.T) {
	values := []string{"a", "b", "c"}
	s := NewSeriesInt("my_series", values)

	if s.Name() != "my_series" {
		t.Errorf("expected name 'my_series', got %s", s.Name())
	}
}

func TestSetName(t *testing.T) {
	values := []int{1, 2, 3}
	s := NewSeriesInt("old_name", values)
	s.SetName("new_name")

	if s.Name() != "new_name" {
		t.Errorf("expected name 'new_name', got %s", s.Name())
	}
}

func TestValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	s := NewSeriesInt("test", values)

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
		s := NewSeriesInt("test", values)

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

		label, value := s.Get(1)
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
		s := NewSeriesInt("test", values)
		s.Get(-1)
	})

	t.Run("panics on out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds index")
			}
		}()
		values := []int{1, 2, 3}
		s := NewSeriesInt("test", values)
		s.Get(10)
	})
}

func TestSeries_IndexGet(t *testing.T) {
	values := []int{1, 2, 3}
	index := []string{"a", "b", "c"}
	s := NewSeries("test", values, index)

	val := s.IndexGet("b")
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
		s := NewSeriesInt("long_series", values)

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
		s := NewSeriesInt("test", values)

		head := s.Head(10)
		if head.Len() != 3 {
			t.Errorf("expected length 3, got %d", head.Len())
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		s := NewSeriesInt("my_series", values)

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
		s := NewSeriesInt("test", values)

		tail := s.Tail(10)
		if tail.Len() != 3 {
			t.Errorf("expected length 3, got %d", tail.Len())
		}
	})

	t.Run("preserves series name", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		s := NewSeriesInt("my_series", values)

		tail := s.Tail(2)
		if tail.Name() != "my_series" {
			t.Errorf("expected name 'my_series', got %s", tail.Name())
		}
	})
}

func TestAppend(t *testing.T) {
	t.Run("appends Series", func(t *testing.T) {
		values := []int{1, 2, 3}
		s := NewSeriesInt("series1", values)

		values = []int{4, 5, 6}
		o := NewSeriesInt("series2", values)

		s.Append(o)
		if s.Len() != 6 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.values[3] != 4 {
			t.Errorf("expected value to be 4, got %d", s.values[3])
		}
	})
}
