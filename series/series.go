package series

import (
	"fmt"
	"strings"
)

// A Series is basically an List which holds Data of a certain type but has extra capabilities
type Series struct {
	name   string
	values []any
	index  []string
}

// creating an new Series
func NewSeries(name string, values []any, index []string) *Series {
	if len(values) == 0 {
		panic("cannot create Series with no data")
	}

	return &Series{
		name:   name,
		values: values,
		index: index,
	}
}

// Len return the length of the value slice
func (s *Series) Len() int {
	return len(s.values)
}

// Name return the name of the Series
func (s *Series) Name() string {
	return s.name
}

// SetSeries sets the Name od the Series
func (s *Series) SetName(name string) {
	s.name = name
}

// Values() returns the values of the series as an copy of the values slice
// through the copy manipulation the original data is impossible
func (s *Series) Values() []any {
	copied := make([]any, len(s.values))
	copy(copied, s.values)
	return copied
}

func (s *Series) Index() []string {
	if s.index == nil {
		index := make([]string, s.Len())
		for i := range index {
			index[i] = fmt.Sprintf("%d", i)
		}
		return index
	}
	return s.index
}

func (s *Series) Get(i int) (string, any) {
	if i < 0 || i >= s.Len() {
		panic(fmt.Sprintf("index %d out of bounds", i))
	}
	label := s.Index()[i] // uses Index() so it generates if nil
	return label, s.values[i]
}

func (s *Series) String() string {
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
		sb.WriteString(fmt.Sprintf("%s: %v\n", label, value))
	}

	if s.Len() > 10 {
		sb.WriteString(fmt.Sprintf("... (%d more)\n", s.Len()-10))
	}

	return sb.String()
}

func (s *Series) Head(n int) *Series { 
	maxlength := n
	if maxlength > s.Len() {
		maxlength = s.Len()
	}
	name := s.name
	values := make([]any, maxlength)
	index := make([]string, maxlength)
	for i := 0; i < maxlength; i++ {
		values[i] = s.values[i]
		index[i] = s.index[i]
	}

	return NewSeries(name, values, index)
}


func (s *Series) Tail(n int) *Series {
	length := s.Len()
	maxlength := n
	if maxlength > length {
		maxlength = length
	}

	name := s.name
	values := make([]any, maxlength)
	index := make([]string, maxlength)
	
	start := length - maxlength
  for i := 0; i < maxlength; i++ {
      values[i] = s.values[start+i]
      index[i] = s.index[start+i]
  }

	return NewSeries(name, values, index)
}
