package series

import "math"

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type NumericSeries[T Numeric, R comparable] struct {
	*Series[T, R]
}

func NewNumericSeries[T Numeric, R comparable](name string, values []T, index []R) *NumericSeries[T, R] {
	return &NumericSeries[T, R]{
		Series: NewSeries(name, values, index),
	}
}

func NewIndexNumericSeries[T Numeric](name string, values []T) *NumericSeries[T, int] {
	index := make([]int, len(values))
	for i := range values {
		index[i] = i
	}

	return NewNumericSeries[T, int](name, values, index)
}

// Sum returns the sum of the Series
func (ns *NumericSeries[T, R]) Sum() T {
	var sum T
	for _, v := range ns.values {
		sum += v
	}
	return sum
}

func (ns *NumericSeries[T, R]) Mean() float64 {
	return float64(ns.Sum()) / float64(ns.Len())
}

// Min returns the smallest value in the Series
func (ns *NumericSeries[T, R]) Min() T {
	if ns.Len() == 0 {
		panic("cannot get min of empty series")
	}

	minValue := ns.values[0]
	for i := 1; i < ns.Len(); i++ {
		if ns.values[i] < minValue {
			minValue = ns.values[i]
		}
	}
	return minValue
}

// Max returns the largest value in the Series
func (ns *NumericSeries[T, R]) Max() T {
	if ns.Len() == 0 {
		panic("cannot get max of empty series")
	}

	maxValue := ns.values[0]
	for i := 1; i < ns.Len(); i++ {
		if ns.values[i] > maxValue {
			maxValue = ns.values[i]
		}
	}
	return maxValue
}

// StdDev returns the standard deviation of the Series
// dof is degrees of freedom, typically 0 for population(complete set) and 1 for sample(uncomplete set)
func (ns *NumericSeries[T, R]) StdDev(dof int) float64 {
	if ns.Len() == 0 {
		panic("cannot get standard deviation of empty series")
	}

	if dof < 0 {
		panic("degrees of freedom must be non-negative")
	}

	mean := ns.Mean()
	var sumSquaredDiff float64

	for _, v := range ns.values {
		diff := float64(v) - mean
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(ns.Len()-dof)
	return math.Sqrt(variance)
}

// Abs returns the Series with absolute values
func (ns *NumericSeries[T, R]) Abs() *NumericSeries[T, R] {
	absValues := make([]T, ns.Len())
	for i, v := range ns.values {
		//math.abs only works for float64
		// and i dont want to convert for because of performance and float arithmetic issues
		if v < 0 {
			absValues[i] = -v
		} else {
			absValues[i] = v
		}
	}
	return NewNumericSeries[T, R](ns.name, absValues, ns.index)
}

// Add adds two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Add(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	if ns.Len() != other.Len() {
		panic("series must be of the same length to add")
	}

	if name == "" {
		name = ns.name + "_add_" + other.name
	}

	resultValues := make([]T, ns.Len())
	for i := range ns.values {
		resultValues[i] = ns.values[i] + other.values[i]
	}
	return NewNumericSeries[T, R](name, resultValues, ns.index)
}

// Subtract subtracts two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Subtract(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	if ns.Len() != other.Len() {
		panic("series must be of the same length to subtract")
	}

	if name == "" {
		name = ns.name + "_sub_" + other.name
	}

	resultValues := make([]T, ns.Len())
	for i := range ns.values {
		resultValues[i] = ns.values[i] - other.values[i]
	}
	return NewNumericSeries[T, R](name, resultValues, ns.index)
}

// ArgMax returns the index position of the largest value in the Series
func (ns *NumericSeries[T, R]) ArgMax() R {
	if ns.Len() == 0 {
		panic("cannot get argmax of empty series")
	}

	maxIndex := 0
	maxValue := ns.values[0]

	for i := 1; i < ns.Len(); i++ {
		if ns.values[i] > maxValue {
			maxValue = ns.values[i]
			maxIndex = i
		}
	}

	return ns.index[maxIndex]
}

// ArgMin returns the index position (int) of the smallest value in the values slice
func (ns *NumericSeries[T, R]) ArgMin() int {
	if ns.Len() == 0 {
		panic("cannot get argmin of empty series")
	}

	minIndex := 0
	minValue := ns.values[0]

	for i := 1; i < ns.Len(); i++ {
		if ns.values[i] < minValue {
			minValue = ns.values[i]
			minIndex = i
		}
	}

	return minIndex
}

// IdxMin returns the label of the smallest value in the Series
func (ns *NumericSeries[T, R]) IdxMin() R {
	return ns.index[ns.ArgMin()]
}

// CumSum returns a new Series with cumulative sum of values
func (ns *NumericSeries[T, R]) CumSum() *NumericSeries[T, R] {
	cumSumValues := make([]T, ns.Len())
	var sum T

	for i, v := range ns.values {
		sum += v
		cumSumValues[i] = sum
	}

	return NewNumericSeries[T, R](ns.name+"_cumsum", cumSumValues, ns.index)
}

// DropNA returns a new Series with NaN values removed (float32 and float64 only)
func (ns *NumericSeries[T, R]) DropNA() *NumericSeries[T, R] {
	var validValues []T
	var validIndex []R

	for i, v := range ns.values {
		// Check if the value is NaN (only for float types)
		isNaN := false
		switch any(v).(type) {
		case float32:
			isNaN = math.IsNaN(float64(any(v).(float32)))
		case float64:
			isNaN = math.IsNaN(any(v).(float64))
		}

		if !isNaN {
			validValues = append(validValues, v)
			validIndex = append(validIndex, ns.index[i])
		}
	}

	if len(validValues) == 0 {
		panic("cannot create series with no data after dropping NA values")
	}

	return NewNumericSeries[T, R](ns.name, validValues, validIndex)
}
