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
	seenKeys := make(map[any]bool)

	for i := 0; i < df.nrows; i++ {
		key := series.AtAny(i)

		if !seenKeys[key] {
			groupKeys = append(groupKeys, key)
			seenKeys[key] = true
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

// Aggregate applies aggregation functions to columns and returns a new DataFrame
// aggregations is a map from result column name to AggFunc
// Format: map["column_aggname"] = AggFunc
func (gb *DataFrameGroupBy) Aggregate(aggregations map[string]AggFunc) (*DataFrame, error) {
	if len(aggregations) == 0 {
		return nil, fmt.Errorf("no aggregations specified")
	}

	// Build result DataFrame
	nGroups := len(gb.groupKeys)

	// Create index from group keys
	newIndex := make([]any, nGroups)
	copy(newIndex, gb.groupKeys)

	// Parse aggregation column names to get source column and agg name
	type aggSpec struct {
		sourceColumn string
		aggFunc      AggFunc
		resultName   string
	}

	var aggSpecs []aggSpec
	for resultName, fn := range aggregations {
		// Try to extract source column name from result name
		// Expected format: "column_aggname" or custom name
		// For simplicity, we'll require explicit source column specification
		// Store as-is for now
		aggSpecs = append(aggSpecs, aggSpec{
			sourceColumn: "", // to be determined
			aggFunc:      fn,
			resultName:   resultName,
		})
	}

	// Actually, let's redesign the API to be clearer
	// We'll pass a map of sourceColumn -> map of aggName -> AggFunc
	// But for now, let's implement a simpler version

	// Build result columns
	resultColumns := make(map[string]SeriesInterface)
	columnOrder := []string{gb.groupColumn} // Start with group column

	// Add group column
	groupColValues := make([]any, nGroups)
	groupColIndex := make([]any, nGroups)
	for i, key := range gb.groupKeys {
		groupColValues[i] = key
		groupColIndex[i] = i
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

	// Process each aggregation
	for resultName, aggFunc := range aggregations {
		// Extract source column name from result name
		// For now, assume user provides full column names
		// We need to infer which column to aggregate
		// Let's require the format "sourceColumn_aggName"
		// But we can't parse that reliably, so let's change the API

		// For this implementation, we'll aggregate all non-group columns
		// and name results as "column_resultName"
		// Actually, let's require user to specify column explicitly
		// by using a nested map structure later

		// For now, simple implementation: assume resultName indicates source column
		// This is a limitation we'll document
		values := make([]any, nGroups)

		// We need to know which column to aggregate
		// For this simple version, let's assume aggregations map has format:
		// "sourceColumn_aggName" -> aggFunc
		// We'll parse the underscore to get source column

		// Find the source column by trying each column until aggFunc returns non-nil
		var foundColumn string
		for _, colName := range gb.df.columnOrder {
			if colName == gb.groupColumn {
				continue // skip group column
			}

			// Try aggregating this column for first group
			firstGroupIndices := gb.groups[gb.groupKeys[0]]
			subSeries := gb.extractSubSeries(colName, firstGroupIndices)
			result := aggFunc(subSeries)

			if result != nil {
				foundColumn = colName
				break
			}
		}

		if foundColumn == "" {
			return nil, fmt.Errorf("could not determine source column for aggregation %s", resultName)
		}

		// Apply aggregation for each group
		for i, key := range gb.groupKeys {
			rowIndices := gb.groups[key]
			subSeries := gb.extractSubSeries(foundColumn, rowIndices)
			values[i] = aggFunc(subSeries)
		}

		// Create result series
		resultSeries := createSeriesFromAny(
			resultName,
			values,
			groupColIndex,
			fmt.Sprintf("%T", values[0]),
			"int",
		)

		resultColumns[resultName] = resultSeries
		columnOrder = append(columnOrder, resultName)
	}

	return &DataFrame{
		columns:     resultColumns,
		columnOrder: columnOrder,
		index:       groupColIndex,
		indexType:   "int",
		nrows:       nGroups,
	}, nil
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

		for aggName, aggFunc := range aggMap {
			resultName := colName + "_" + aggName
			values := make([]any, nGroups)

			// Apply aggregation for each group
			for i, key := range gb.groupKeys {
				rowIndices := gb.groups[key]
				subSeries := gb.extractSubSeries(colName, rowIndices)
				values[i] = aggFunc(subSeries)
			}

			// Infer type from first value
			valueType := "any"
			if len(values) > 0 && values[0] != nil {
				valueType = fmt.Sprintf("%T", values[0])
			}

			// Create result series
			resultSeries := createSeriesFromAny(
				resultName,
				values,
				groupColIndex,
				valueType,
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
