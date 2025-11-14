package series

import (
	"fmt"
	"math"
	"reflect"
)

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

// Operation executes a custom operation on two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Operation(other *NumericSeries[T, R], op func(a, b T) T, name string) *NumericSeries[T, R] {
	if ns.Len() != other.Len() {
		panic("series must be of the same length to perform operation")
	}

	if name == "" {
		name = ns.name + "_op_" + other.name
	}

	resultValues := make([]T, ns.Len())
	for i := range ns.values {
		resultValues[i] = op(ns.values[i], other.values[i])
	}
	return NewNumericSeries[T, R](name, resultValues, ns.index)
}

// Add adds two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Add(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	addfunc := func(a, b T) T {
		return a + b
	}
	return ns.Operation(other, addfunc, name)
}

// Subtract subtracts two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Subtract(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	subtractfunc := func(a, b T) T {
		return a - b
	}
	return ns.Operation(other, subtractfunc, name)
}

// Multiply multiplies two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Multiply(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	multiplyfunc := func(a, b T) T {
		return a * b
	}
	return ns.Operation(other, multiplyfunc, name)
}

// Divide divides two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Divide(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	dividefunc := func(a, b T) T {
		return a / b
	}
	return ns.Operation(other, dividefunc, name)
}

// Mod does modulus of two NumericSeries element-wise
// use an empty string for name to get the default name
func (ns *NumericSeries[T, R]) Mod(other *NumericSeries[T, R], name string) *NumericSeries[T, R] {
	var zero T
	switch any(zero).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// ok, continue
	default:
		panic("modulus operation only supported for whole number types")
	}

	modfunc := func(a, b T) T {
		// Use reflection to perform modulus
		aVal := reflect.ValueOf(a)
		bVal := reflect.ValueOf(b)

		result := aVal.Int() % bVal.Int()
		return reflect.ValueOf(result).Convert(reflect.TypeOf(a)).Interface().(T)
	}
	return ns.Operation(other, modfunc, name)
}

// Pow raises each element of the Series to the given power
func (ns *NumericSeries[T, R]) Pow(power float64, name string) *NumericSeries[T, R] {
	if name == "" {
		name = fmt.Sprintf("%v_pow(%v)", ns.name, power)
	}

	powValues := make([]T, ns.Len())
	for i, v := range ns.values {
		result := math.Pow(float64(v), power)
		powValues[i] = T(result)
	}

	return NewNumericSeries[T, R](name, powValues, ns.index)
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

// CoVariance computes the covariance between two NumericSeries
// dof is degrees of freedom, typically 0 for population(complete set) and 1 for sample(uncomplete set)
func (ns *NumericSeries[T, R]) CoVariance(other *NumericSeries[T, R], dof int) float64 {
	if ns.Len() != other.Len() {
		panic("series must be of the same length to compute covariance")
	}

	if dof < 0 {
		panic("degrees of freedom must be non-negative")
	}

	meanX := ns.Mean()
	meanY := other.Mean()

	var covSum float64
	for i := 0; i < ns.Len(); i++ {
		covSum += (float64(ns.values[i]) - meanX) * (float64(other.values[i]) - meanY)
	}
	return covSum / float64(ns.Len()-dof)
}

// Correlation computes the Pearson correlation coefficient between two NumericSeries
func (ns *NumericSeries[T, R]) Correlation(other *NumericSeries[T, R]) float64 {
	if ns.Len() != other.Len() {
		panic("series must be of the same length to compute correlation")
	}

	stdX := ns.StdDev(0)
	stdY := other.StdDev(0)
	coVariance := ns.CoVariance(other, 0)

	denominator := stdX * stdY
	if denominator != 0 {
		return coVariance / denominator
	}

	return 0.0
}
