package dataframe

import (
	"fmt"
)

// DataFrameGroupBy represents a grouped DataFrame
type DataFrameGroupBy struct {
	df          *DataFrame
	groupColumn string
	groups      map[any][]int // maps group key to row indices
	groupKeys   []any         // maintains order of groups
}

// GroupBy groups the DataFrame by values in the specified column
func (df *DataFrame) GroupBy(columnName string) (*DataFrameGroupBy, error) {
	series, err := df.GetColumn(columnName)
	if err != nil {
		return nil, err
	}

	// Build groups map
	groups := make(map[any][]int)
	var groupKeys []any

	for i := 0; i < df.nrows; i++ {
		key := series.AtAny(i)

		if _, ok := groups[key]; !ok {
			groupKeys = append(groupKeys, key)
		}

		groups[key] = append(groups[key], i)
	}

	return &DataFrameGroupBy{
		df:          df,
		groupColumn: columnName,
		groups:      groups,
		groupKeys:   groupKeys,
	}, nil
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

// AggregateColumns applies specific aggregations to specific columns
// aggregations is a map from column name to map of result suffix to AggFunc
// Example: {"price": {"mean": AggMean, "sum": AggSum}} creates "price_mean" and "price_sum"
func (gb *DataFrameGroupBy) AggregateColumns(aggregations map[string]map[string]AggFunc) (*DataFrame, error) {
	if len(aggregations) == 0 {
		return nil, fmt.Errorf("no aggregations specified")
	}

	// Build result DataFrame
	nGroups := len(gb.groupKeys)

	// Create index from group keys
	groupColIndex := make([]any, nGroups)
	for i := range nGroups {
		groupColIndex[i] = i
	}

	// Build result columns
	resultColumns := make(map[string]SeriesInterface)
	columnOrder := []string{gb.groupColumn}

	// Add group column
	groupColValues := make([]any, nGroups)
	for i, key := range gb.groupKeys {
		groupColValues[i] = key
	}

	groupSeries, err := gb.df.GetColumn(gb.groupColumn)
	if err != nil {
		return nil, err
	}

	resultColumns[gb.groupColumn] = createSeriesFromAny(
		gb.groupColumn,
		groupColValues,
		groupColIndex,
		groupSeries.GetValueType(),
		"int",
	)

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
		for i, key := range gb.groupKeys {
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
func (gb *DataFrameGroupBy) Count() *DataFrame {
	nGroups := len(gb.groupKeys)

	groupColIndex := make([]any, nGroups)
	for i := range nGroups {
		groupColIndex[i] = i
	}

	// Group column values
	groupColValues := make([]any, nGroups)
	for i, key := range gb.groupKeys {
		groupColValues[i] = key
	}

	groupSeries, _ := gb.df.GetColumn(gb.groupColumn)

	// Count values
	countValues := make([]any, nGroups)
	for i, key := range gb.groupKeys {
		countValues[i] = len(gb.groups[key])
	}

	resultColumns := make(map[string]SeriesInterface)
	resultColumns[gb.groupColumn] = createSeriesFromAny(
		gb.groupColumn,
		groupColValues,
		groupColIndex,
		groupSeries.GetValueType(),
		"int",
	)
	resultColumns["count"] = createSeriesFromAny(
		"count",
		countValues,
		groupColIndex,
		"int",
		"int",
	)

	return &DataFrame{
		columns:     resultColumns,
		columnOrder: []string{gb.groupColumn, "count"},
		index:       groupColIndex,
		indexType:   "int",
		nrows:       nGroups,
	}
}
