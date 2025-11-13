package series

import (
	"fmt"
	"strings"
)

// A Series is basically a List which holds Data of a certain type but has extra capabilities
type Series[T any, R comparable] struct {
	name   string
	values []T
	index  []R
}

// NewSeries creates a new Series
func NewSeries[T any, R comparable](name string, values []T, index []R) *Series[T, R] {
	if len(values) == 0 {
		panic("cannot create Series with no data")
	}

	if index == nil {
		panic("needs an index if you dont have one use IndexedSeries")
	}

	if len(index) != len(values) {
		panic("index length must match values length")
	}

	return &Series[T, R]{
		name:   name,
		values: values,
		index:  index,
	}
}

// Len return the length of the value slice
func (s *Series[T, R]) Len() int {
	return len(s.values)
}

// Name return the name of the Series
func (s *Series[T, R]) Name() string {
	return s.name
}

// SetName sets the Name od the Series
func (s *Series[T, R]) SetName(name string) {
	s.name = name
}

// Values returns the values of the series as an copy of the values slice
// through the copy manipulation the original data is impossible
func (s *Series[T, R]) Values() []T {
	copied := make([]T, len(s.values))
	copy(copied, s.values)
	return copied
}

// Index returns the index of the series as a copy of the index slice
func (s *Series[T, R]) Index() []R {
	copied := make([]R, len(s.index))
	copy(copied, s.index)
	return copied
}

// Get returns the value for the given label
func (s *Series[T, R]) Get(label R) T {
	for i := range s.index {
		if s.index[i] == label {
			return s.values[i]
		}
	}
	panic(fmt.Sprintf("no value found for label %v", label))
}

// GetLabel returns the label and value for the given label
//Not need as you already have your label
/*func (s *Series[T, R]) GetLabel(label R) (R, T) {
	for i := range s.index {
		if s.index[i] == label {
			return s.index[i], s.values[i]
		}
	}
	var zero T
	var zeroR R
	return zeroR, zero
}*/

// At returns the value at the given index od the slice
func (s *Series[T, R]) At(i int) T {
	if i < 0 || i >= s.Len() {
		panic(fmt.Sprintf("index %d out of bounds", i))
	}
	return s.values[i]
}

// AtIndex returns the label and value at the given index of the slice
func (s *Series[T, R]) AtIndex(i int) (R, T) {
	if i < 0 || i >= s.Len() {
		panic(fmt.Sprintf("index %d out of bounds", i))
	}
	return s.index[i], s.values[i]
}

// String returns a string representation of the Series
func (s *Series[T, R]) String() string {
	var sb strings.Builder

	if s.name != "" {
		sb.WriteString(s.name + "\n")
	}

	// don't print more than 10 rows
	maxLen := min(s.Len(), 10)

	for i := range maxLen {
		label, value := s.AtIndex(i)
		sb.WriteString(fmt.Sprintf("%v: %v\n", label, value))
	}

	if s.Len() > 10 {
		sb.WriteString(fmt.Sprintf("... (%d more)\n", s.Len()-10))
	}

	return sb.String()
}

// Head returns the first n elements of the Series
func (s *Series[T, R]) Head(n int) *Series[T, R] {
	maxlength := min(n, s.Len())
	name := s.name
	values := make([]T, maxlength)
	index := make([]R, maxlength)
	for i := range maxlength {
		values[i] = s.values[i]
		index[i] = s.index[i]
	}

	return NewSeries(name, values, index)
}

// Tail returns the last n elements of the Series
func (s *Series[T, R]) Tail(n int) *Series[T, R] {
	length := s.Len()
	maxlength := min(n, length)

	name := s.name
	values := make([]T, maxlength)
	index := make([]R, maxlength)

	start := length - maxlength
	for i := range maxlength {
		values[i] = s.values[start+i]
		index[i] = s.index[start+i]
	}

	return NewSeries(name, values, index)
}

// Append appends another Series to the end of this Series
func (s *Series[T, R]) Append(o *Series[T, R]) {
	s.values = append(s.values, o.values...)
	s.index = append(s.index, o.index...)
}

// Prepend prepends another Series to the beginning of this Series
func (s *Series[T, R]) Prepend(o *Series[T, R]) {
	s.values = append(o.values, s.values...)
	s.index = append(o.index, s.index...)
}

// isEmpty checks if the series is empty
func (s *Series[T, R]) isEmpty() bool {
	return s.Len() == 0
}

// ResetIndex resets the index to default 0..n-1
func (s *Series[T, R]) ResetIndex() *Series[T, int] {
	newIndex := make([]int, s.Len())
	for i := range newIndex {
		newIndex[i] = i
	}
	return NewSeries[T, int](s.name, s.values, newIndex)
}

// SetIndex returns a new Series with the given index
func SetIndex[T any, R comparable, S comparable](s *Series[T, R], newIndex []S) *Series[T, S] {
	if len(newIndex) != s.Len() {
		panic("new index length must match values length")
	}
	return NewSeries[T, S](s.name, s.values, newIndex)
}
