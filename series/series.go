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

type SeriesInt[T any] Series[T, int]

// NewSeries creates an new Series
func NewSeries[T any, R comparable](name string, values []T, index []R) *Series[T, R] {
	if len(values) == 0 {
		panic("cannot create Series with no data")
	}

	if index == nil {
		panic("needs an index if you dont have one use SeriesInt")
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

func NewSeriesInt[T any](name string, values []T) *Series[T, int] {
	index := make([]int, len(values))
	for i := range values {
		index[i] = i
	}

	return NewSeries[T, int](name, values, index)
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

func (s *Series[T, R]) Index() []R {
	copied := make([]R, len(s.index))
	copy(copied, s.index)
	return copied
}

func (s *Series[T, R]) IndexGet(label R) T {
	for i := range s.index {
		if s.index[i] == label {
			return s.values[i]
		}
	}
	var zero T
	return zero
}

func (s *SeriesInt[T]) extendsIndex() []int {
	if s.index == nil {
		s.index = make([]int, len(s.values))
	}
	for i := range s.index {
		if s.index[0] == 0 && i != 0 {
			s.index[i] = i
		}
	}
	return s.index
}

func (s *Series[T, R]) Get(i int) (R, T) {
	if i < 0 || i >= s.Len() {
		panic(fmt.Sprintf("index %d out of bounds", i))
	}
	label := s.index[i]
	return label, s.values[i]
}

func (s *Series[T, R]) String() string {
	var sb strings.Builder

	if s.name != "" {
		sb.WriteString(s.name + "\n")
	}

	maxLen := 10 // don't print more than 10 rows
	if s.Len() < maxLen {
		maxLen = s.Len()
	}

	for i := 0; i < maxLen; i++ {
		label, value := s.Get(i)
		sb.WriteString(fmt.Sprintf("%v: %v\n", label, value))
	}

	if s.Len() > 10 {
		sb.WriteString(fmt.Sprintf("... (%d more)\n", s.Len()-10))
	}

	return sb.String()
}

func (s *Series[T, R]) Head(n int) *Series[T, R] {
	maxlength := n
	if maxlength > s.Len() {
		maxlength = s.Len()
	}
	name := s.name
	values := make([]T, maxlength)
	index := make([]R, maxlength)
	for i := 0; i < maxlength; i++ {
		values[i] = s.values[i]
		index[i] = s.index[i]
	}

	return NewSeries(name, values, index)
}

func (s *Series[T, R]) Tail(n int) *Series[T, R] {
	length := s.Len()
	maxlength := n
	if maxlength > length {
		maxlength = length
	}

	name := s.name
	values := make([]T, maxlength)
	index := make([]R, maxlength)

	start := length - maxlength
	for i := 0; i < maxlength; i++ {
		values[i] = s.values[start+i]
		index[i] = s.index[start+i]
	}

	return NewSeries(name, values, index)
}

func (s *Series[T, R]) Append(o *Series[T, R]) {
	s.values = append(s.values, o.values...)
	s.index = append(s.index, o.index...)
}
