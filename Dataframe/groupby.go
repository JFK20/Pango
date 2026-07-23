package dataframe

import (
	"fmt"
	"strings"
)

// DataFrameGroupBy represents a grouped DataFrame
type DataFrameGroupBy struct {
	df           *DataFrame
	groupColumns []string
	groups       map[string][]int // maps composite group key to row indices
	groupOrder   []string         // maintains order of groups
	keyValues    map[string][]any // composite key -> typed tuple, parallel to groupColumns
}

// compositeKey joins a row's grouped-column values into a single hashable map key
func compositeKey(tuple []any) string {
	parts := make([]string, len(tuple))
	for i, v := range tuple {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(parts, "\x1f")
}

// GroupBy groups the DataFrame by values in the specified column
func (df *DataFrame) GroupBy(columnName string) (*DataFrameGroupBy, error) {
	return df.GroupByColumns(columnName)
}

// GroupByColumns groups the DataFrame by the composite values of the specified columns
func (df *DataFrame) GroupByColumns(columnNames ...string) (*DataFrameGroupBy, error) {
	if len(columnNames) == 0 {
		return nil, fmt.Errorf("must specify at least one column")
	}

	cols := make([]SeriesInterface, len(columnNames))
	for i, name := range columnNames {
		s, err := df.GetColumn(name)
		if err != nil {
			return nil, err
		}
		cols[i] = s
	}

	groups := make(map[string][]int)
	keyValues := make(map[string][]any)
	var groupOrder []string

	for i := 0; i < df.nrows; i++ {
		tuple := make([]any, len(cols))
		for j, c := range cols {
			tuple[j] = c.AtAny(i)
		}
		key := compositeKey(tuple)

		if _, ok := groups[key]; !ok {
			groupOrder = append(groupOrder, key)
			keyValues[key] = tuple
		}

		groups[key] = append(groups[key], i)
	}

	orderedColumns := make([]string, len(columnNames))
	copy(orderedColumns, columnNames)

	return &DataFrameGroupBy{
		df:           df,
		groupColumns: orderedColumns,
		groups:       groups,
		groupOrder:   groupOrder,
		keyValues:    keyValues,
	}, nil
}

// Groups returns a map from each group's key to the DataFrame containing that group's rows.
// For a single-column groupby the key is the grouped column's raw value. For a
// GroupByColumns groupby the key is an internal composite identifier - recover the
// individual grouped values from the corresponding columns in the returned DataFrame.
func (gb *DataFrameGroupBy) Groups() map[any]*DataFrame {
	result := make(map[any]*DataFrame, len(gb.groupOrder))
	for _, key := range gb.groupOrder {
		result[gb.publicKey(key)] = gb.df.selectRows(gb.groups[key])
	}
	return result
}

// Range provides an iterator over each group's key and DataFrame, in the order groups
// were first encountered.
func (gb *DataFrameGroupBy) Range() func(yield func(key any, group *DataFrame) bool) {
	return func(yield func(key any, group *DataFrame) bool) {
		for _, key := range gb.groupOrder {
			if !yield(gb.publicKey(key), gb.df.selectRows(gb.groups[key])) {
				return
			}
		}
	}
}

// publicKey returns the key exposed to callers of Groups/Range for a composite key
func (gb *DataFrameGroupBy) publicKey(key string) any {
	if len(gb.groupColumns) == 1 {
		return gb.keyValues[key][0]
	}
	return key
}

// AggFunc is a function that aggregates a SeriesInterface into a single value
type AggFunc func(SeriesInterface) any

// Predefined aggregation functions

// AggSum returns the sum of a numeric series
func AggSum(s SeriesInterface) any {
	if ns, ok := s.(NumericSeriesInterface); ok {
		return ns.SumFloat()
	}
	return nil
}

// AggMean returns the mean of a numeric series
func AggMean(s SeriesInterface) any {
	if ns, ok := s.(NumericSeriesInterface); ok {
		return ns.Mean()
	}
	return nil
}

// AggMin returns the minimum of a numeric series
func AggMin(s SeriesInterface) any {
	if ns, ok := s.(NumericSeriesInterface); ok {
		return ns.MinFloat()
	}
	return nil
}

// AggMax returns the maximum of a numeric series
func AggMax(s SeriesInterface) any {
	if ns, ok := s.(NumericSeriesInterface); ok {
		return ns.MaxFloat()
	}
	return nil
}

// AggCount returns the count of elements
func AggCount(s SeriesInterface) any {
	return s.Len()
}

// AggFirst returns the first value in the series
func AggFirst(s SeriesInterface) any {
	if s.Len() > 0 {
		return s.AtAny(0)
	}
	return nil
}

// AggLast returns the last value in the series
func AggLast(s SeriesInterface) any {
	if s.Len() > 0 {
		return s.AtAny(s.Len() - 1)
	}
	return nil
}

// AggStdDev returns the standard deviation of a numeric series (sample)
func AggStdDev(s SeriesInterface) any {
	if ns, ok := s.(NumericSeriesInterface); ok {
		return ns.StdDev(1) // use sample standard deviation (dof=1)
	}
	return nil
}

// AggQuantile returns an AggFunc computing the given quantile (0-1) of a numeric series
func AggQuantile(q float64) AggFunc {
	return func(s SeriesInterface) any {
		if ns, ok := s.(NumericSeriesInterface); ok {
			return ns.Quantile(q)
		}
		return nil
	}
}

// AggMedian returns an AggFunc computing the median of a numeric series
func AggMedian() AggFunc {
	return func(s SeriesInterface) any {
		if ns, ok := s.(NumericSeriesInterface); ok {
			return ns.Median()
		}
		return nil
	}
}

// AggregateColumns applies specific aggregations to specific columns
// aggregations is a map from column name to map of result suffix to AggFunc
// Example: {"price": {"mean": AggMean, "sum": AggSum}} creates "price_mean" and "price_sum"
func (gb *DataFrameGroupBy) AggregateColumns(aggregations map[string]map[string]AggFunc) (*DataFrame, error) {
	if len(aggregations) == 0 {
		return nil, fmt.Errorf("no aggregations specified")
	}

	// Build result DataFrame
	nGroups := len(gb.groupOrder)

	// Create index from group keys
	groupColIndex := make([]any, nGroups)
	for i := range nGroups {
		groupColIndex[i] = i
	}

	// Build result columns
	resultColumns := make(map[string]SeriesInterface)
	columnOrder := make([]string, len(gb.groupColumns))
	copy(columnOrder, gb.groupColumns)

	// Add one result column per grouped column
	for gi, groupColName := range gb.groupColumns {
		groupSeries, err := gb.df.GetColumn(groupColName)
		if err != nil {
			return nil, err
		}

		groupColValues := make([]any, nGroups)
		for i, key := range gb.groupOrder {
			groupColValues[i] = gb.keyValues[key][gi]
		}

		resultColumns[groupColName] = createSeriesFromAny(
			groupColName,
			groupColValues,
			groupColIndex,
			groupSeries.GetValueType(),
			"int",
		)
	}

	// Process each column and its aggregations
	for colName, aggMap := range aggregations {
		// Verify column exists
		if _, err := gb.df.GetColumn(colName); err != nil {
			return nil, fmt.Errorf("column %s not found: %w", colName, err)
		}

		// Build each group's sub-series once, shared across every
		// aggregation requested for this column (previously rebuilt once
		// per aggregation, i.e. once per group per agg function).
		subSeriesByGroup := make([]SeriesInterface, nGroups)
		for i, key := range gb.groupOrder {
			subSeriesByGroup[i] = gb.extractSubSeries(colName, gb.groups[key])
		}

		for aggName, aggFunc := range aggMap {
			resultName := colName + "_" + aggName
			values := make([]any, nGroups)

			for i, subSeries := range subSeriesByGroup {
				values[i] = aggFunc(subSeries)
			}

			// Create result series
			resultSeries := createSeriesFromAny(
				resultName,
				values,
				groupColIndex,
				inferValueType(values),
				"int",
			)

			resultColumns[resultName] = resultSeries
			columnOrder = append(columnOrder, resultName)
		}
	}

	return &DataFrame{
		columns:     resultColumns,
		columnOrder: columnOrder,
		index:       groupColIndex,
		indexType:   "int",
		nrows:       nGroups,
	}, nil
}

// extractSubSeries creates a SeriesInterface containing only the rows specified by indices
func (gb *DataFrameGroupBy) extractSubSeries(columnName string, rowIndices []int) SeriesInterface {
	series := gb.df.columns[columnName]

	values := make([]any, len(rowIndices))
	index := make([]any, len(rowIndices))

	for i, rowIdx := range rowIndices {
		idx, val := series.AtIndexAny(rowIdx)
		index[i] = idx
		values[i] = val
	}

	return createSeriesFromAny(
		columnName,
		values,
		index,
		series.GetValueType(),
		series.GetIndexType(),
	)
}

// Count returns the size of each group
func (gb *DataFrameGroupBy) Count() (*DataFrame, error) {
	nGroups := len(gb.groupOrder)

	groupColIndex := make([]any, nGroups)
	for i := range nGroups {
		groupColIndex[i] = i
	}

	resultColumns := make(map[string]SeriesInterface)
	columnOrder := make([]string, len(gb.groupColumns))
	copy(columnOrder, gb.groupColumns)

	for gi, groupColName := range gb.groupColumns {
		groupSeries, err := gb.df.GetColumn(groupColName)
		if err != nil {
			return nil, err
		}

		groupColValues := make([]any, nGroups)
		for i, key := range gb.groupOrder {
			groupColValues[i] = gb.keyValues[key][gi]
		}

		resultColumns[groupColName] = createSeriesFromAny(
			groupColName,
			groupColValues,
			groupColIndex,
			groupSeries.GetValueType(),
			"int",
		)
	}

	// Count values
	countValues := make([]any, nGroups)
	for i, key := range gb.groupOrder {
		countValues[i] = len(gb.groups[key])
	}

	resultColumns["count"] = createSeriesFromAny(
		"count",
		countValues,
		groupColIndex,
		"int",
		"int",
	)

	columnOrder = append(columnOrder, "count")

	return &DataFrame{
		columns:     resultColumns,
		columnOrder: columnOrder,
		index:       groupColIndex,
		indexType:   "int",
		nrows:       nGroups,
	}, nil
}
