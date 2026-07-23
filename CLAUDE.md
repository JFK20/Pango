# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Pango is a Go library bringing pandas-like data manipulation to Go (name = "pandas" + "Go"), built on Go generics for type-safe series/dataframe operations. It is in early development. The stated project goal (README.md) is to port pandas' most important features while keeping 100% test coverage — the `series` package currently holds to this (`go test ./series/... -cover` ≈ 99%); `Dataframe` has real but lower coverage (`go test ./Dataframe/... -coverpkg=pango/Dataframe -cover` ≈ 93%), so treat changes there with extra care and prefer adding tests alongside new code.

## Commands

```bash
go build ./...              # build everything
go run main.go               # run the example programs in main.go
go test ./...                 # run all tests
go test ./series/...          # run tests for one package
go test ./series/... -run TestName   # run a single test by name
go test ./series/... -cover   # check coverage (keep at/near 100% for series)
go vet ./...                  # static checks
```

There is no linter config beyond `go vet`; there are no build tags or Makefile.

## Architecture

Two packages form the core, plus an `example/` tree with runnable usage examples driven from `main.go`:

- **`series/`** — the foundational data type.
  - `Series[T comparable, R comparable]` (`series.go`) is a generic labeled 1-D array: `T` is the value type, `R` is the index/label type. Constructed via `NewSeries`, or `NewIndexSeries`/`IndexedSeries` (`indexSeries.go`) for a default `0..n-1` int index.
  - `NumericSeries[T Numeric, R comparable]` (`seriesNumeric.go`) embeds `*Series[T, R]` and adds numeric ops (`Sum`, `Mean`, `Min`/`Max`, `StdDev`, `Quantile`/`Median`, `CoVariance`, `Correlation`, element-wise `Add`/`Subtract`/`Multiply`/`Divide`/`Mod`/`Pow`, `CumSum`, `DropNA`, `ArgMin`/`ArgMax`, etc.). `Numeric` is a constraint covering all int/uint/float kinds. `Quantile`/`Median` (and `Min`/`Max`/`StdDev`/`ArgMin`/`ArgMax`) panic on an empty series or invalid arguments rather than returning an error, consistent with the rest of the file.
  - Most mutating-looking operations (`Head`, `Tail`, `Copy`, `SortByIndex`, `SortByValue`, etc.) return **new** series rather than mutating in place — `Values()`/`Index()` return defensive copies.
  - Every type also implements an `Any`-suffixed shadow API (`AtAny`, `ValuesAny`, `CopyAny`, `GetValueType`, `GetIndexType`, etc.) purely so it can satisfy `dataframe.SeriesInterface` / `dataframe.NumericSeriesInterface` — Go generics can't be used across the `Dataframe` package boundary directly, so this boxing-to-`any` layer is how heterogeneous typed columns get stored in a single `DataFrame`.

- **`Dataframe/`** — tabular structure built on top of `series` via interfaces, not concrete types.
  - `seriesInterface.go` defines `SeriesInterface` (what any column must support) and `NumericSeriesInterface` (adds `SumFloat`/`Mean`/`MinFloat`/`MaxFloat`/`Count`/`StdDev`/`Quantile`/`Median`, required for numeric aggregation). `Dataframe` never imports `series` directly — it only depends on these interfaces, which `series.Series`/`series.NumericSeries` satisfy structurally. The `Dataframe`-internal `genericNumericSeries` (see below) is the interface's other implementer, so any addition to `NumericSeriesInterface` needs a matching method there too or the package stops compiling.
  - `dataframe.go` holds `DataFrame` (`columns map[string]SeriesInterface` + `columnOrder []string` + `index`/`indexType`/`nrows`) and its transforms (`Select`, `DropColumn`, `Head`/`Tail`, `FilterColumn`/`FilterIn`, `ApplyToColumn(s)`, `AstypeColumn`, `ResetIndex`, `SetIndexFromColumn`, etc.). `NewDataFrame`/`NewDataFrameWithIndex` store independent copies of each input series via `CopyAny()`, so a caller mutating its own typed `Series` after construction can't desync the `DataFrame`. Any transform that rebuilds a column from raw `[]any` (e.g. after filtering rows) goes through `createSeriesFromAny`, which picks a numeric-capable wrapper (`genericNumericSeries`) when the boxed values are all numeric, or a plain `genericSeries` otherwise — this is what lets a derived/filtered/grouped column still support `GetNumericColumn`. Row-subsetting operations (`FilterColumn`, and `DataFrameGroupBy.Groups`/`Range`) share the private `selectRows(indices []int) *DataFrame` helper.
  - `groupby.go` implements `GroupBy`/`GroupByColumns`/`DataFrameGroupBy` and aggregation via `AggFunc` (`SeriesInterface -> any`). `GroupByColumns(columnNames ...string)` supports composite multi-column grouping: rows are keyed internally by a `\x1f`-joined string built from each grouped column's stringified value (`compositeKey`), since a `[]any` tuple can't be used as a Go map key; the original typed values are kept alongside in `keyValues` so aggregation output preserves types. `AggregateColumns(map[colName]map[aggSuffix]AggFunc)` is the row-collapsing aggregation entry point (result columns are named `col_aggsuffix`); `Count()` is the same shape but just counts. Built-in `AggFunc`s (`AggSum`, `AggMean`, `AggMin`, `AggMax`, `AggStdDev`, `AggQuantile`, `AggMedian`, `AggCount`, `AggFirst`, `AggLast`) type-assert to `NumericSeriesInterface` where relevant; `AggQuantile(q)`/`AggMedian()` are factories that return an `AggFunc`, not `AggFunc`s themselves. For raw (non-collapsed) per-group access, `Groups() map[any]*DataFrame` and `Range() func(yield func(key any, group *DataFrame) bool)` hand back each group's rows as its own `*DataFrame` (which still contains the grouped-by columns) — the exposed key is the raw scalar for a single-column `GroupBy`, or the opaque composite string for `GroupByColumns`.

`main.go` just calls into two example packages, each with three example flows: `example/series/examples.go` (`Example1SumColumn`, `Example2AddColumns`, `Example3MultipleOperations`) shows the intended `series`-only usage pattern, retrieving typed numeric columns via a type assertion to the concrete `*series.NumericSeries[T, R]` when type-specific ops (e.g. `Add`/`Multiply`) are needed; `example/dataframes/examples.go` (`Example1SelectFilterHeadTail`, `Example2TransformColumns`, `Example3GroupByAggregate`) shows the `DataFrame` side — building a `DataFrame` from `series`, `Select`/`FilterColumn`/`Head`/`Tail`, `ApplyToColumn`, and `GroupBy` + `AggregateColumns`.
