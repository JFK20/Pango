package dataframe

// SeriesInterface defines the common interface for all Series types
// This allows DataFrame to work with different typed Series in a type-safe way
type SeriesInterface interface {
	// Len returns the number of elements in the Series
	Len() int

	// Name returns the name of the Series
	Name() string

	// SetName sets the name of the Series
	SetName(name string)

	// AtAny returns the value at the given position as any
	AtAny(i int) any

	// AtIndexAny returns both the index label and value at the given position as any
	AtIndexAny(i int) (any, any)

	// CopyAny creates a deep copy of the Series
	CopyAny() SeriesInterface

	// IndexAny returns the index labels as []any
	IndexAny() []any

	// ValuesAny returns the values as []any
	ValuesAny() []any

	// GetValueType returns the type of values stored in the Series as a string
	GetValueType() string

	// GetIndexType returns the type of index stored in the Series as a string
	GetIndexType() string
}

// NumericSeriesInterface extends SeriesInterface with numeric operations
// This is used for aggregation functions that require numeric data
type NumericSeriesInterface interface {
	SeriesInterface

	// SumFloat returns the sum of all values as float64
	SumFloat() float64

	// Mean returns the arithmetic mean of all values
	Mean() float64

	// MinFloat returns the minimum value as float64
	MinFloat() float64

	// MaxFloat returns the maximum value as float64
	MaxFloat() float64

	// Count returns the number of elements (same as Len but included for aggregation consistency)
	Count() int

	// StdDev returns the standard deviation with given degrees of freedom
	StdDev(dof int) float64

	// Quantile returns the value at the given quantile (0-1) using linear interpolation
	Quantile(q float64) float64

	// Median returns the median value (Quantile(0.5))
	Median() float64
}
