package series

import (
	"math"
	"testing"
)

func TestNewNumericSeries(t *testing.T) {
	t.Run("creates numeric series with index", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		ns := NewNumericSeries("test", values, index)

		if ns.Len() != 3 {
			t.Errorf("expected length 3, got %d", ns.Len())
		}
		if ns.Name() != "test" {
			t.Errorf("expected name 'test', got %s", ns.Name())
		}
	})
}

func TestNewIndexNumericSeries(t *testing.T) {
	t.Run("creates numeric series with default index", func(t *testing.T) {
		values := []float64{1.5, 2.5, 3.5}
		ns := NewIndexNumericSeries("test", values)

		if ns.Len() != 3 {
			t.Errorf("expected length 3, got %d", ns.Len())
		}

		idx := ns.Index()
		for i := range idx {
			if idx[i] != i {
				t.Errorf("expected index %d at position %d, got %d", i, i, idx[i])
			}
		}
	})
}

func TestSum(t *testing.T) {
	t.Run("sums integer values", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		ns := NewIndexNumericSeries("test", values)

		sum := ns.Sum()
		if sum != 15 {
			t.Errorf("expected sum 15, got %d", sum)
		}
	})

	t.Run("sums float values", func(t *testing.T) {
		values := []float64{1.5, 2.5, 3.5}
		ns := NewIndexNumericSeries("test", values)

		sum := ns.Sum()
		if sum != 7.5 {
			t.Errorf("expected sum 7.5, got %f", sum)
		}
	})

	t.Run("sums negative values", func(t *testing.T) {
		values := []int{-1, -2, -3}
		ns := NewIndexNumericSeries("test", values)

		sum := ns.Sum()
		if sum != -6 {
			t.Errorf("expected sum -6, got %d", sum)
		}
	})
}

func TestMean(t *testing.T) {
	t.Run("calculates mean of integer values", func(t *testing.T) {
		values := []int{2, 4, 6, 8, 10}
		ns := NewIndexNumericSeries("test", values)

		mean := ns.Mean()
		if mean != 6.0 {
			t.Errorf("expected mean 6.0, got %f", mean)
		}
	})

	t.Run("calculates mean of float values", func(t *testing.T) {
		values := []float64{1.0, 2.0, 3.0}
		ns := NewIndexNumericSeries("test", values)

		mean := ns.Mean()
		if mean != 2.0 {
			t.Errorf("expected mean 2.0, got %f", mean)
		}
	})
}

func TestMin(t *testing.T) {
	t.Run("finds minimum integer value", func(t *testing.T) {
		values := []int{5, 2, 8, 1, 9}
		ns := NewIndexNumericSeries("test", values)

		minimum := ns.Min()
		if minimum != 1 {
			t.Errorf("expected min 1, got %d", minimum)
		}
	})

	t.Run("finds minimum float value", func(t *testing.T) {
		values := []float64{5.5, 2.2, 8.8, 1.1, 9.9}
		ns := NewIndexNumericSeries("test", values)

		minimum := ns.Min()
		if minimum != 1.1 {
			t.Errorf("expected min 1.1, got %f", minimum)
		}
	})

	// is this necessary this isnt javaScript
	t.Run("finds minimum with negative values", func(t *testing.T) {
		values := []int{5, -10, 8, 1, 9}
		ns := NewIndexNumericSeries("test", values)

		minimum := ns.Min()
		if minimum != -10 {
			t.Errorf("expected min -10, got %d", minimum)
		}
	})

	t.Run("panics on empty series", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty series")
			}
		}()

		values := []int{0}
		ns := NewIndexNumericSeries("test", values)
		ns.values = []int{}
		ns.Min()
	})
}

func TestMax(t *testing.T) {
	t.Run("finds maximum integer value", func(t *testing.T) {
		values := []int{5, 2, 8, 1, 9}
		ns := NewIndexNumericSeries("test", values)

		maximum := ns.Max()
		if maximum != 9 {
			t.Errorf("expected max 9, got %d", maximum)
		}
	})

	t.Run("finds maximum float value", func(t *testing.T) {
		values := []float64{5.5, 2.2, 8.8, 1.1, 9.9}
		ns := NewIndexNumericSeries("test", values)

		maximum := ns.Max()
		if maximum != 9.9 {
			t.Errorf("expected max 9.9, got %f", maximum)
		}
	})

	t.Run("finds maximum with negative values", func(t *testing.T) {
		values := []int{-5, -10, -8, -1, -9}
		ns := NewIndexNumericSeries("test", values)

		maximum := ns.Max()
		if maximum != -1 {
			t.Errorf("expected max -1, got %d", maximum)
		}
	})

	t.Run("panics on empty series", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty series")
			}
		}()

		values := []int{0}
		ns := NewIndexNumericSeries("test", values)
		ns.values = []int{}
		ns.Max()
	})
}

func TestStdDev(t *testing.T) {
	t.Run("calculates standard deviation with dof 0", func(t *testing.T) {
		values := []float64{2.0, 4.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0}
		ns := NewIndexNumericSeries("test", values)

		stdDev := ns.StdDev(0)
		expected := 1.9116
		if math.Abs(stdDev-expected) > 0.01 {
			t.Errorf("expected standard deviation ~%f, got %f", expected, stdDev)
		}
	})

	t.Run("calculates standard deviation with dof 1", func(t *testing.T) {
		values := []float64{2.0, 4.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0}
		ns := NewIndexNumericSeries("test", values)

		stdDev := ns.StdDev(1)
		expected := 2.0276
		if math.Abs(stdDev-expected) > 0.01 {
			t.Errorf("expected standard deviation ~%f, got %f", expected, stdDev)
		}
	})

	t.Run("standard deviation of uniform values is zero", func(t *testing.T) {
		values := []int{5, 5, 5, 5, 5}
		ns := NewIndexNumericSeries("test", values)

		stdDev := ns.StdDev(0)
		if stdDev != 0.0 {
			t.Errorf("expected standard deviation 0, got %f", stdDev)
		}

		stdDev = ns.StdDev(1)
		if stdDev != 0.0 {
			t.Errorf("expected standard deviation 0, got %f", stdDev)
		}
	})

	t.Run("panics on negative dof", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative dof")
			}
		}()

		values := []int{1, 2, 3}
		ns := NewIndexNumericSeries("test", values)
		ns.StdDev(-1)
	})

	t.Run("panics on empty series", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty series")
			}
		}()

		values := []int{0}
		ns := NewIndexNumericSeries("test", values)
		ns.values = []int{}
		ns.StdDev(0)
	})
}

func TestAbs(t *testing.T) {
	t.Run("returns absolute values for integers", func(t *testing.T) {
		values := []int{-5, 2, -8, 1, -9}
		ns := NewIndexNumericSeries("test", values)

		absNs := ns.Abs()
		expected := []int{5, 2, 8, 1, 9}

		for i := range expected {
			if absNs.At(i) != expected[i] {
				t.Errorf("expected %d at position %d, got %d", expected[i], i, absNs.At(i))
			}
		}
	})

	t.Run("returns absolute values for floats", func(t *testing.T) {
		values := []float64{-5.5, 2.2, -8.8, 1.1, -9.9}
		ns := NewIndexNumericSeries("test", values)

		absNs := ns.Abs()
		expected := []float64{5.5, 2.2, 8.8, 1.1, 9.9}

		for i := range expected {
			if absNs.At(i) != expected[i] {
				t.Errorf("expected %f at position %d, got %f", expected[i], i, absNs.At(i))
			}
		}
	})

	t.Run("preserves positive values", func(t *testing.T) {
		values := []int{5, 2, 8, 1, 9}
		ns := NewIndexNumericSeries("test", values)

		absNs := ns.Abs()

		for i := range values {
			if absNs.At(i) != values[i] {
				t.Errorf("expected %d at position %d, got %d", values[i], i, absNs.At(i))
			}
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("adds two integer series", func(t *testing.T) {
		values1 := []int{1, 2, 3}
		values2 := []int{4, 5, 6}
		ns1 := NewIndexNumericSeries("series1", values1)
		ns2 := NewIndexNumericSeries("series2", values2)

		result := ns1.Add(ns2, "")
		expected := []int{5, 7, 9}

		for i := range expected {
			if result.At(i) != expected[i] {
				t.Errorf("expected %d at position %d, got %d", expected[i], i, result.At(i))
			}
		}
	})

	t.Run("adds two float series", func(t *testing.T) {
		values1 := []float64{1.5, 2.5, 3.5}
		values2 := []float64{0.5, 0.5, 0.5}
		ns1 := NewIndexNumericSeries("series1", values1)
		ns2 := NewIndexNumericSeries("series2", values2)

		result := ns1.Add(ns2, "result")
		expected := []float64{2.0, 3.0, 4.0}

		if result.Name() != "result" {
			t.Errorf("expected name 'result', got %s", result.Name())
		}

		for i := range expected {
			if result.At(i) != expected[i] {
				t.Errorf("expected %f at position %d, got %f", expected[i], i, result.At(i))
			}
		}
	})

	t.Run("uses default name when empty string provided", func(t *testing.T) {
		values1 := []int{1, 2, 3}
		values2 := []int{4, 5, 6}
		ns1 := NewIndexNumericSeries("s1", values1)
		ns2 := NewIndexNumericSeries("s2", values2)

		result := ns1.Add(ns2, "")
		if result.Name() != "s1_add_s2" {
			t.Errorf("expected default name 's1_add_s2', got %s", result.Name())
		}
	})

	t.Run("panics with different lengths", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for different lengths")
			}
		}()

		values1 := []int{1, 2, 3}
		values2 := []int{4, 5}
		ns1 := NewIndexNumericSeries("series1", values1)
		ns2 := NewIndexNumericSeries("series2", values2)

		ns1.Add(ns2, "")
	})
}

func TestSubtract(t *testing.T) {
	t.Run("subtracts two integer series", func(t *testing.T) {
		values1 := []int{10, 20, 30}
		values2 := []int{1, 2, 3}
		ns1 := NewIndexNumericSeries("series1", values1)
		ns2 := NewIndexNumericSeries("series2", values2)

		result := ns1.Subtract(ns2, "")
		expected := []int{9, 18, 27}

		for i := range expected {
			if result.At(i) != expected[i] {
				t.Errorf("expected %d at position %d, got %d", expected[i], i, result.At(i))
			}
		}
	})

	t.Run("subtracts two float series", func(t *testing.T) {
		values1 := []float64{5.5, 10.5, 15.5}
		values2 := []float64{1.5, 2.5, 3.5}
		ns1 := NewIndexNumericSeries("series1", values1)
		ns2 := NewIndexNumericSeries("series2", values2)

		result := ns1.Subtract(ns2, "difference")
		expected := []float64{4.0, 8.0, 12.0}

		if result.Name() != "difference" {
			t.Errorf("expected name 'difference', got %s", result.Name())
		}

		for i := range expected {
			if result.At(i) != expected[i] {
				t.Errorf("expected %f at position %d, got %f", expected[i], i, result.At(i))
			}
		}
	})

	t.Run("uses default name when empty string provided", func(t *testing.T) {
		values1 := []int{10, 20, 30}
		values2 := []int{1, 2, 3}
		ns1 := NewIndexNumericSeries("s1", values1)
		ns2 := NewIndexNumericSeries("s2", values2)

		result := ns1.Subtract(ns2, "")
		if result.Name() != "s1_sub_s2" {
			t.Errorf("expected default name 's1_sub_s2', got %s", result.Name())
		}
	})

	t.Run("panics with different lengths", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for different lengths")
			}
		}()

		values1 := []int{10, 20, 30}
		values2 := []int{1, 2}
		ns1 := NewIndexNumericSeries("series1", values1)
		ns2 := NewIndexNumericSeries("series2", values2)

		ns1.Subtract(ns2, "")
	})
}

func TestArgMax(t *testing.T) {
	t.Run("returns label of maximum value", func(t *testing.T) {
		values := []int{5, 2, 8, 1, 9}
		index := []string{"a", "b", "c", "d", "e"}
		ns := NewNumericSeries("test", values, index)

		argMax := ns.ArgMax()
		if argMax != "e" {
			t.Errorf("expected argmax 'e', got %s", argMax)
		}
	})

	t.Run("returns first occurrence when multiple maxima", func(t *testing.T) {
		values := []int{5, 9, 8, 9, 1}
		index := []string{"a", "b", "c", "d", "e"}
		ns := NewNumericSeries("test", values, index)

		argMax := ns.ArgMax()
		if argMax != "b" {
			t.Errorf("expected argmax 'b', got %s", argMax)
		}
	})

	t.Run("panics on empty series", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty series")
			}
		}()

		values := []int{0}
		ns := NewIndexNumericSeries("test", values)
		ns.values = []int{}
		ns.ArgMax()
	})
}

func TestArgMin(t *testing.T) {
	t.Run("returns index position of minimum value", func(t *testing.T) {
		values := []int{5, 2, 8, 1, 9}
		ns := NewIndexNumericSeries("test", values)

		argMin := ns.ArgMin()
		if argMin != 3 {
			t.Errorf("expected argmin 3, got %d", argMin)
		}
	})

	t.Run("returns first occurrence when multiple minima", func(t *testing.T) {
		values := []int{5, 1, 8, 1, 9}
		ns := NewIndexNumericSeries("test", values)

		argMin := ns.ArgMin()
		if argMin != 1 {
			t.Errorf("expected argmin 1, got %d", argMin)
		}
	})

	t.Run("panics on empty series", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty series")
			}
		}()

		values := []int{0}
		ns := NewIndexNumericSeries("test", values)
		ns.values = []int{}
		ns.ArgMin()
	})
}

func TestIdxMin(t *testing.T) {
	t.Run("returns label of minimum value", func(t *testing.T) {
		values := []int{5, 2, 8, 1, 9}
		index := []string{"a", "b", "c", "d", "e"}
		ns := NewNumericSeries("test", values, index)

		idxMin := ns.IdxMin()
		if idxMin != "d" {
			t.Errorf("expected idxmin 'd', got %s", idxMin)
		}
	})

	t.Run("works with integer index", func(t *testing.T) {
		values := []float64{5.5, 2.2, 8.8, 1.1, 9.9}
		ns := NewIndexNumericSeries("test", values)

		idxMin := ns.IdxMin()
		if idxMin != 3 {
			t.Errorf("expected idxmin 3, got %d", idxMin)
		}
	})
}

func TestCumSum(t *testing.T) {
	t.Run("calculates cumulative sum for integers", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		ns := NewIndexNumericSeries("test", values)

		cumSum := ns.CumSum()
		expected := []int{1, 3, 6, 10, 15}

		if cumSum.Name() != "test_cumsum" {
			t.Errorf("expected name 'test_cumsum', got %s", cumSum.Name())
		}

		for i := range expected {
			if cumSum.At(i) != expected[i] {
				t.Errorf("expected %d at position %d, got %d", expected[i], i, cumSum.At(i))
			}
		}
	})

	t.Run("calculates cumulative sum for floats", func(t *testing.T) {
		values := []float64{1.5, 2.5, 3.5}
		ns := NewIndexNumericSeries("test", values)

		cumSum := ns.CumSum()
		expected := []float64{1.5, 4.0, 7.5}

		for i := range expected {
			if cumSum.At(i) != expected[i] {
				t.Errorf("expected %f at position %d, got %f", expected[i], i, cumSum.At(i))
			}
		}
	})

	t.Run("preserves index", func(t *testing.T) {
		values := []int{1, 2, 3}
		index := []string{"a", "b", "c"}
		ns := NewNumericSeries("test", values, index)

		cumSum := ns.CumSum()
		resultIndex := cumSum.Index()

		for i := range index {
			if resultIndex[i] != index[i] {
				t.Errorf("expected index %s at position %d, got %s", index[i], i, resultIndex[i])
			}
		}
	})
}

func TestDropNA(t *testing.T) {
	t.Run("removes NaN values from float64 series", func(t *testing.T) {
		values := []float64{1.0, math.NaN(), 3.0, math.NaN(), 5.0}
		ns := NewIndexNumericSeries("test", values)

		cleaned := ns.DropNA()
		if cleaned.Len() != 3 {
			t.Errorf("expected length 3, got %d", cleaned.Len())
		}

		expected := []float64{1.0, 3.0, 5.0}
		for i := range expected {
			if cleaned.At(i) != expected[i] {
				t.Errorf("expected %f at position %d, got %f", expected[i], i, cleaned.At(i))
			}
		}
	})

	t.Run("removes NaN values from float32 series", func(t *testing.T) {
		values := []float32{1.0, float32(math.NaN()), 3.0, 4.0}
		ns := NewIndexNumericSeries("test", values)

		cleaned := ns.DropNA()
		if cleaned.Len() != 3 {
			t.Errorf("expected length 3, got %d", cleaned.Len())
		}
	})

	t.Run("preserves all values for integer series", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		ns := NewIndexNumericSeries("test", values)

		cleaned := ns.DropNA()
		if cleaned.Len() != 5 {
			t.Errorf("expected length 5, got %d", cleaned.Len())
		}
	})

	t.Run("preserves corresponding indices", func(t *testing.T) {
		values := []float64{1.0, math.NaN(), 3.0, math.NaN(), 5.0}
		index := []string{"a", "b", "c", "d", "e"}
		ns := NewNumericSeries("test", values, index)

		cleaned := ns.DropNA()
		resultIndex := cleaned.Index()

		expectedIndex := []string{"a", "c", "e"}
		for i := range expectedIndex {
			if resultIndex[i] != expectedIndex[i] {
				t.Errorf("expected index %s at position %d, got %s", expectedIndex[i], i, resultIndex[i])
			}
		}
	})

	t.Run("panics when all values are NaN", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when all values are NaN")
			}
		}()

		values := []float64{math.NaN(), math.NaN(), math.NaN()}
		ns := NewIndexNumericSeries("test", values)
		ns.DropNA()
	})

	t.Run("returns unchanged series when no NaN values", func(t *testing.T) {
		values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		ns := NewIndexNumericSeries("test", values)

		cleaned := ns.DropNA()
		if cleaned.Len() != 5 {
			t.Errorf("expected length 5, got %d", cleaned.Len())
		}

		for i := range values {
			if cleaned.At(i) != values[i] {
				t.Errorf("expected %f at position %d, got %f", values[i], i, cleaned.At(i))
			}
		}
	})
}
