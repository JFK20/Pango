package dataframe

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// DataFrame is a tabular data structure with labeled columns
// Each column is a SeriesInterface allowing heterogeneous typed columns
type DataFrame struct {
	columns     map[string]SeriesInterface
	columnOrder []string
	index       []any
	indexType   string
	nrows       int
}

// NewDataFrame creates a new DataFrame from a variadic list of Series
// It uses the index from the first Series or creates a default 0..n-1 index
func NewDataFrame(series ...SeriesInterface) (*DataFrame, error) {
	if len(series) == 0 {
		return nil, fmt.Errorf("cannot create DataFrame with no series")
	}

	// Check that all series have the same length
	nrows := series[0].Len()
	for i, s := range series {
		if s.Len() != nrows {
			return nil, fmt.Errorf("series %d (%s) has length %d, expected %d", i, s.Name(), s.Len(), nrows)
		}
	}

	// Check for duplicate column names
	namesSeen := make(map[string]bool)
	for _, s := range series {
		name := s.Name()
		if namesSeen[name] {
			return nil, fmt.Errorf("duplicate column name: %s", name)
		}
		namesSeen[name] = true
	}

	// Use index from first series
	index := series[0].IndexAny()
	indexType := series[0].GetIndexType()

	// Build columns map and order
	columns := make(map[string]SeriesInterface)
	columnOrder := make([]string, len(series))
	for i, s := range series {
		name := s.Name()
		columns[name] = s
		columnOrder[i] = name
	}

	return &DataFrame{
		columns:     columns,
		columnOrder: columnOrder,
		index:       index,
		indexType:   indexType,
		nrows:       nrows,
	}, nil
}

// NewDataFrameWithIndex creates a new DataFrame with a custom index
func NewDataFrameWithIndex(index []any, series ...SeriesInterface) (*DataFrame, error) {
	if len(series) == 0 {
		return nil, fmt.Errorf("cannot create DataFrame with no series")
	}

	nrows := len(index)

	// Check that all series have the same length as index
	for i, s := range series {
		if s.Len() != nrows {
			return nil, fmt.Errorf("series %d (%s) has length %d, expected %d", i, s.Name(), s.Len(), nrows)
		}
	}

	// Check for duplicate column names
	namesSeen := make(map[string]bool)
	for _, s := range series {
		name := s.Name()
		if namesSeen[name] {
			return nil, fmt.Errorf("duplicate column name: %s", name)
		}
		namesSeen[name] = true
	}

	// Infer index type from first element
	indexType := "any"
	if len(index) > 0 && index[0] != nil {
		indexType = reflect.TypeOf(index[0]).String()
	}

	// Build columns map and order
	columns := make(map[string]SeriesInterface)
	columnOrder := make([]string, len(series))
	for i, s := range series {
		name := s.Name()
		columns[name] = s
		columnOrder[i] = name
	}

	return &DataFrame{
		columns:     columns,
		columnOrder: columnOrder,
		index:       index,
		indexType:   indexType,
		nrows:       nrows,
	}, nil
}

// Len returns the number of rows in the DataFrame
func (df *DataFrame) Len() int {
	return df.nrows
}

// Columns returns the column names in order
func (df *DataFrame) Columns() []string {
	result := make([]string, len(df.columnOrder))
	copy(result, df.columnOrder)
	return result
}

// Shape returns the dimensions (rows, columns) of the DataFrame
func (df *DataFrame) Shape() (int, int) {
	return df.nrows, len(df.columnOrder)
}

// GetColumn retrieves a column by name
func (df *DataFrame) GetColumn(name string) (SeriesInterface, error) {
	s, ok := df.columns[name]
	if !ok {
		return nil, fmt.Errorf("column %s not found", name)
	}
	return s, nil
}

// GetNumericColumn retrieves a column as NumericSeriesInterface
func (df *DataFrame) GetNumericColumn(name string) (NumericSeriesInterface, error) {
	s, err := df.GetColumn(name)
	if err != nil {
		return nil, err
	}

	ns, ok := s.(NumericSeriesInterface)
	if !ok {
		return nil, fmt.Errorf("column %s is not a numeric series (type: %s)", name, s.GetValueType())
	}
	return ns, nil
}

// AddColumn adds a new column to the DataFrame
func (df *DataFrame) AddColumn(s SeriesInterface) error {
	if s.Len() != df.nrows {
		return fmt.Errorf("series length %d does not match DataFrame length %d", s.Len(), df.nrows)
	}

	name := s.Name()
	if _, exists := df.columns[name]; exists {
		return fmt.Errorf("column %s already exists", name)
	}

	df.columns[name] = s
	df.columnOrder = append(df.columnOrder, name)
	return nil
}

// DropColumn removes a column and returns a new DataFrame
func (df *DataFrame) DropColumn(name string) (*DataFrame, error) {
	if _, ok := df.columns[name]; !ok {
		return nil, fmt.Errorf("column %s not found", name)
	}

	// Create new columns map and order
	newColumns := make(map[string]SeriesInterface)
	newOrder := make([]string, 0, len(df.columnOrder)-1)

	for _, colName := range df.columnOrder {
		if colName != name {
			newColumns[colName] = df.columns[colName].CopyAny()
			newOrder = append(newOrder, colName)
		}
	}

	// Copy index
	newIndex := make([]any, len(df.index))
	copy(newIndex, df.index)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       df.nrows,
	}, nil
}

// RenameColumn renames a column
func (df *DataFrame) RenameColumn(oldName, newName string) error {
	if _, ok := df.columns[oldName]; !ok {
		return fmt.Errorf("column %s not found", oldName)
	}

	if _, ok := df.columns[newName]; ok {
		return fmt.Errorf("column %s already exists", newName)
	}

	// Update map. Rename on a copy so we don't mutate a series object the
	// caller may still hold a reference to (columns are stored by reference,
	// not copied, when a DataFrame is constructed).
	s := df.columns[oldName].CopyAny()
	s.SetName(newName)
	df.columns[newName] = s
	delete(df.columns, oldName)

	// Update order
	for i, name := range df.columnOrder {
		if name == oldName {
			df.columnOrder[i] = newName
			break
		}
	}

	return nil
}

// Index returns a copy of the DataFrame index
func (df *DataFrame) Index() []any {
	result := make([]any, len(df.index))
	copy(result, df.index)
	return result
}

// String returns a string representation of the DataFrame
func (df *DataFrame) String() string {
	var sb strings.Builder

	rows, cols := df.Shape()
	sb.WriteString(fmt.Sprintf("DataFrame [%d rows x %d columns]\n", rows, cols))

	if rows == 0 || cols == 0 {
		return sb.String()
	}

	// Header
	sb.WriteString(fmt.Sprintf("%-10s", "Index"))
	for _, colName := range df.columnOrder {
		sb.WriteString(fmt.Sprintf("%-15s", colName))
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString(strings.Repeat("-", 10+15*cols))
	sb.WriteString("\n")

	// Rows (max 10)
	maxRows := min(rows, 10)
	for i := 0; i < maxRows; i++ {
		sb.WriteString(fmt.Sprintf("%-10v", df.index[i]))
		for _, colName := range df.columnOrder {
			val := df.columns[colName].AtAny(i)
			sb.WriteString(fmt.Sprintf("%-15v", val))
		}
		sb.WriteString("\n")
	}

	if rows > 10 {
		sb.WriteString(fmt.Sprintf("... (%d more rows)\n", rows-10))
	}

	return sb.String()
}

// ColumnTypes returns a map of column names to their value types
func (df *DataFrame) ColumnTypes() map[string]string {
	result := make(map[string]string)
	for name, series := range df.columns {
		result[name] = series.GetValueType()
	}
	return result
}

// Info returns a string with information about the DataFrame
func (df *DataFrame) Info() string {
	var sb strings.Builder

	rows, cols := df.Shape()
	sb.WriteString(fmt.Sprintf("DataFrame Info:\n"))
	sb.WriteString(fmt.Sprintf("  Rows: %d\n", rows))
	sb.WriteString(fmt.Sprintf("  Columns: %d\n", cols))
	sb.WriteString(fmt.Sprintf("  Index Type: %s\n\n", df.indexType))

	sb.WriteString("Column Details:\n")
	for _, colName := range df.columnOrder {
		series := df.columns[colName]
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", colName, series.GetValueType()))
	}

	return sb.String()
}

// Select returns a new DataFrame with only the specified columns
func (df *DataFrame) Select(columnNames ...string) (*DataFrame, error) {
	if len(columnNames) == 0 {
		return nil, fmt.Errorf("must specify at least one column")
	}

	// Verify all columns exist
	for _, name := range columnNames {
		if _, ok := df.columns[name]; !ok {
			return nil, fmt.Errorf("column %s not found", name)
		}
	}

	// Create new columns with deep copies
	newColumns := make(map[string]SeriesInterface)
	newOrder := make([]string, len(columnNames))

	for i, name := range columnNames {
		newColumns[name] = df.columns[name].CopyAny()
		newOrder[i] = name
	}

	// Copy index
	newIndex := make([]any, len(df.index))
	copy(newIndex, df.index)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       df.nrows,
	}, nil
}

// Head returns a new DataFrame with the first n rows
func (df *DataFrame) Head(n int) *DataFrame {
	if n <= 0 {
		n = 5 // default to 5 rows
	}

	maxLen := min(n, df.nrows)

	// Copy index
	newIndex := make([]any, maxLen)
	copy(newIndex, df.index[:maxLen])

	// Create subset of each series
	newColumns := make(map[string]SeriesInterface)
	for name, series := range df.columns {
		newSeries := createSubseries(series, 0, maxLen)
		newColumns[name] = newSeries
	}

	// Copy column order
	newOrder := make([]string, len(df.columnOrder))
	copy(newOrder, df.columnOrder)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       maxLen,
	}
}

// Tail returns a new DataFrame with the last n rows
func (df *DataFrame) Tail(n int) *DataFrame {
	if n <= 0 {
		n = 5 // default to 5 rows
	}

	maxLen := min(n, df.nrows)
	start := df.nrows - maxLen

	// Copy index
	newIndex := make([]any, maxLen)
	copy(newIndex, df.index[start:])

	// Create subset of each series
	newColumns := make(map[string]SeriesInterface)
	for name, series := range df.columns {
		newSeries := createSubseries(series, start, df.nrows)
		newColumns[name] = newSeries
	}

	// Copy column order
	newOrder := make([]string, len(df.columnOrder))
	copy(newOrder, df.columnOrder)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       maxLen,
	}
}

// createSubseries creates a new series with a subset of rows
// This is a helper function for Head/Tail/Filter operations
func createSubseries(s SeriesInterface, start, end int) SeriesInterface {
	// Build new values and index slices
	length := end - start
	newValues := make([]any, length)
	newIndex := make([]any, length)

	for i := 0; i < length; i++ {
		idx, val := s.AtIndexAny(start + i)
		newIndex[i] = idx
		newValues[i] = val
	}

	// We need to create a typed Series from []any
	// Use reflection to determine types and create appropriate Series
	return createSeriesFromAny(s.Name(), newValues, newIndex, s.GetValueType(), s.GetIndexType())
}

// createSeriesFromAny creates a SeriesInterface from []any slices.
// If every non-nil value is numeric, it returns a genericNumericSeries so
// that the result still satisfies NumericSeriesInterface (needed for Head,
// Tail, FilterColumn, ApplyToColumn and every groupby aggregation, which all
// route their per-row/per-group data through this function).
func createSeriesFromAny(name string, values []any, index []any, valueType, indexType string) SeriesInterface {
	base := genericSeries{
		name:      name,
		values:    values,
		index:     index,
		valueType: valueType,
		indexType: indexType,
	}
	if isNumericValues(values) {
		return &genericNumericSeries{genericSeries: base}
	}
	return &base
}

// toFloat64 converts a boxed numeric value to float64.
func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

// inferValueType returns the Go type name of the first non-nil value, or
// "any" if every value is nil.
func inferValueType(values []any) string {
	for _, v := range values {
		if v != nil {
			return fmt.Sprintf("%T", v)
		}
	}
	return "any"
}

// isNumericValues reports whether every non-nil value is numeric, and there
// is at least one such value.
func isNumericValues(values []any) bool {
	seenNumeric := false
	for _, v := range values {
		if v == nil {
			continue
		}
		if _, ok := toFloat64(v); !ok {
			return false
		}
		seenNumeric = true
	}
	return seenNumeric
}

// genericSeries is a fallback implementation of SeriesInterface
// Used when we can't determine the concrete type
type genericSeries struct {
	name      string
	values    []any
	index     []any
	valueType string
	indexType string
}

func (g *genericSeries) Len() int                    { return len(g.values) }
func (g *genericSeries) Name() string                { return g.name }
func (g *genericSeries) SetName(name string)         { g.name = name }
func (g *genericSeries) AtAny(i int) any             { return g.values[i] }
func (g *genericSeries) AtIndexAny(i int) (any, any) { return g.index[i], g.values[i] }
func (g *genericSeries) IndexAny() []any             { return g.index }
func (g *genericSeries) ValuesAny() []any            { return g.values }
func (g *genericSeries) GetValueType() string        { return g.valueType }
func (g *genericSeries) GetIndexType() string        { return g.indexType }
func (g *genericSeries) CopyAny() SeriesInterface {
	newValues := make([]any, len(g.values))
	copy(newValues, g.values)
	newIndex := make([]any, len(g.index))
	copy(newIndex, g.index)
	return &genericSeries{
		name:      g.name,
		values:    newValues,
		index:     newIndex,
		valueType: g.valueType,
		indexType: g.indexType,
	}
}

// genericNumericSeries extends genericSeries with NumericSeriesInterface
// support by converting its boxed values to float64 on demand.
type genericNumericSeries struct {
	genericSeries
}

func (g *genericNumericSeries) floatValues() []float64 {
	result := make([]float64, 0, len(g.values))
	for _, v := range g.values {
		if f, ok := toFloat64(v); ok {
			result = append(result, f)
		}
	}
	return result
}

func (g *genericNumericSeries) SumFloat() float64 {
	var sum float64
	for _, f := range g.floatValues() {
		sum += f
	}
	return sum
}

func (g *genericNumericSeries) Mean() float64 {
	fv := g.floatValues()
	if len(fv) == 0 {
		return 0
	}
	var sum float64
	for _, f := range fv {
		sum += f
	}
	return sum / float64(len(fv))
}

func (g *genericNumericSeries) MinFloat() float64 {
	fv := g.floatValues()
	if len(fv) == 0 {
		return 0
	}
	min := fv[0]
	for _, f := range fv[1:] {
		if f < min {
			min = f
		}
	}
	return min
}

func (g *genericNumericSeries) MaxFloat() float64 {
	fv := g.floatValues()
	if len(fv) == 0 {
		return 0
	}
	max := fv[0]
	for _, f := range fv[1:] {
		if f > max {
			max = f
		}
	}
	return max
}

func (g *genericNumericSeries) Count() int { return g.Len() }

func (g *genericNumericSeries) StdDev(dof int) float64 {
	fv := g.floatValues()
	n := len(fv)
	if n == 0 || n-dof <= 0 {
		return 0
	}
	mean := g.Mean()
	var sumSquaredDiff float64
	for _, f := range fv {
		diff := f - mean
		sumSquaredDiff += diff * diff
	}
	return math.Sqrt(sumSquaredDiff / float64(n-dof))
}

func (g *genericNumericSeries) CopyAny() SeriesInterface {
	copied := g.genericSeries.CopyAny().(*genericSeries)
	return &genericNumericSeries{genericSeries: *copied}
}

// Range provides an iterator over rows of the DataFrame
// Returns a function that yields index and row map for each row
func (df *DataFrame) Range() func(yield func(idx any, row map[string]any) bool) {
	return func(yield func(idx any, row map[string]any) bool) {
		for i := 0; i < df.nrows; i++ {
			row := make(map[string]any)
			for _, colName := range df.columnOrder {
				row[colName] = df.columns[colName].AtAny(i)
			}
			if !yield(df.index[i], row) {
				return
			}
		}
	}
}

// FilterColumn returns a new DataFrame with rows where the predicate is true for the specified column
func (df *DataFrame) FilterColumn(columnName string, predicate func(any) bool) (*DataFrame, error) {
	series, err := df.GetColumn(columnName)
	if err != nil {
		return nil, err
	}

	// Collect indices of rows that match predicate
	var matchingIndices []int
	for i := 0; i < df.nrows; i++ {
		if predicate(series.AtAny(i)) {
			matchingIndices = append(matchingIndices, i)
		}
	}

	// Build new index
	newIndex := make([]any, len(matchingIndices))
	for i, idx := range matchingIndices {
		newIndex[i] = df.index[idx]
	}

	// Build new columns with filtered data
	newColumns := make(map[string]SeriesInterface)
	for _, colName := range df.columnOrder {
		col := df.columns[colName]
		values := make([]any, len(matchingIndices))
		for i, idx := range matchingIndices {
			values[i] = col.AtAny(idx)
		}
		newColumns[colName] = createSeriesFromAny(colName, values, newIndex, col.GetValueType(), col.GetIndexType())
	}

	// Copy column order
	newOrder := make([]string, len(df.columnOrder))
	copy(newOrder, df.columnOrder)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       len(matchingIndices),
	}, nil
}

// ResetIndex returns a new DataFrame with index reset to 0..n-1
func (df *DataFrame) ResetIndex() *DataFrame {
	newIndex := make([]any, df.nrows)
	for i := range df.nrows {
		newIndex[i] = i
	}

	// Deep copy columns
	newColumns := make(map[string]SeriesInterface)
	for name, series := range df.columns {
		newColumns[name] = series.CopyAny()
	}

	// Copy column order
	newOrder := make([]string, len(df.columnOrder))
	copy(newOrder, df.columnOrder)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   "int",
		nrows:       df.nrows,
	}
}

// SetIndexFromColumn sets the index from a column's values and removes that column
func (df *DataFrame) SetIndexFromColumn(columnName string) (*DataFrame, error) {
	series, err := df.GetColumn(columnName)
	if err != nil {
		return nil, err
	}

	// Get new index from column values
	newIndex := series.ValuesAny()
	indexType := series.GetValueType()

	// Create new columns excluding the index column
	newColumns := make(map[string]SeriesInterface)
	newOrder := make([]string, 0, len(df.columnOrder)-1)

	for _, colName := range df.columnOrder {
		if colName != columnName {
			newColumns[colName] = df.columns[colName].CopyAny()
			newOrder = append(newOrder, colName)
		}
	}

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   indexType,
		nrows:       df.nrows,
	}, nil
}

// ApplyToColumn applies a function to each element of a column and returns a new DataFrame
// The function receives and returns any type, allowing flexible transformations
func (df *DataFrame) ApplyToColumn(columnName string, fn func(any) any) (*DataFrame, error) {
	series, err := df.GetColumn(columnName)
	if err != nil {
		return nil, err
	}

	// Apply function to each value
	newValues := make([]any, df.nrows)
	for i := 0; i < df.nrows; i++ {
		newValues[i] = fn(series.AtAny(i))
	}

	// Create new series with transformed values
	newIndex := make([]any, len(df.index))
	copy(newIndex, df.index)

	newSeries := createSeriesFromAny(columnName, newValues, newIndex, inferValueType(newValues), df.indexType)

	// Create new DataFrame with replaced column
	newColumns := make(map[string]SeriesInterface)
	for _, colName := range df.columnOrder {
		if colName == columnName {
			newColumns[colName] = newSeries
		} else {
			newColumns[colName] = df.columns[colName].CopyAny()
		}
	}

	// Copy column order
	newOrder := make([]string, len(df.columnOrder))
	copy(newOrder, df.columnOrder)

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       df.nrows,
	}, nil
}

// ApplyToColumns applies a function across multiple columns element-wise
// The function receives a map of column values for each row
func (df *DataFrame) ApplyToColumns(columns []string, fn func(map[string]any) any, resultColumnName string) (*DataFrame, error) {
	// Verify all columns exist
	for _, colName := range columns {
		if _, ok := df.columns[colName]; !ok {
			return nil, fmt.Errorf("column %s not found", colName)
		}
	}

	// Apply function row by row
	newValues := make([]any, df.nrows)
	for i := 0; i < df.nrows; i++ {
		row := make(map[string]any)
		for _, colName := range columns {
			row[colName] = df.columns[colName].AtAny(i)
		}
		newValues[i] = fn(row)
	}

	// Create new series
	newIndex := make([]any, len(df.index))
	copy(newIndex, df.index)

	newSeries := createSeriesFromAny(resultColumnName, newValues, newIndex, inferValueType(newValues), df.indexType)

	// Create new DataFrame with added column
	newColumns := make(map[string]SeriesInterface)
	for name, series := range df.columns {
		newColumns[name] = series.CopyAny()
	}
	newColumns[resultColumnName] = newSeries

	// Update column order
	newOrder := make([]string, len(df.columnOrder)+1)
	copy(newOrder, df.columnOrder)
	newOrder[len(df.columnOrder)] = resultColumnName

	return &DataFrame{
		columns:     newColumns,
		columnOrder: newOrder,
		index:       newIndex,
		indexType:   df.indexType,
		nrows:       df.nrows,
	}, nil
}
