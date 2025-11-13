package series

// IndexedSeries is a Series with integer index starting from 0
type IndexedSeries[T any] Series[T, int]

// NewIndexSeries creates a new Series with integer index starting from 0
func NewIndexSeries[T any](name string, values []T) *Series[T, int] {
	index := make([]int, len(values))
	for i := range values {
		index[i] = i
	}

	return NewSeries(name, values, index)
}
