# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Pango is a Go library bringing pandas-like data manipulation to Go (name = "pandas" + "Go"), built on Go generics for type-safe series/dataframe operations. It is in early development. The stated project goal (README.md) is to port pandas' most important features while keeping 100% test coverage — the `series` package currently holds to this (`go test ./series/... -cover` ≈ 99%); `Dataframe` has no tests yet, so treat changes there with extra care and prefer adding tests alongside new code.

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

Two packages form the core, plus a `main.go` with runnable usage examples:

- **`series/`** — the foundational data type.
  - `Series[T comparable, R comparable]` (`series.go`) is a generic labeled 1-D array: `T` is the value type, `R` is the index/label type. Constructed via `NewSeries`, or `NewIndexSeries`/`IndexedSeries` (`indexSeries.go`) for a default `0..n-1` int index.
  - `NumericSeries[T Numeric, R comparable]` (`seriesNumeric.go`) embeds `*Series[T, R]` and adds numeric ops (`Sum`, `Mean`, `Min`/`Max`, `StdDev`, `CoVariance`, `Correlation`, element-wise `Add`/`Subtract`/`Multiply`/`Divide`/`Mod`/`Pow`, `CumSum`, `DropNA`, `ArgMin`/`ArgMax`, etc.). `Numeric` is a constraint covering all int/uint/float kinds.
  - Most mutating-looking operations (`Head`, `Tail`, `Copy`, `SortByIndex`, `SortByValue`, etc.) return **new** series rather than mutating in place — `Values()`/`Index()` return defensive copies.
  - Every type also implements an `Any`-suffixed shadow API (`AtAny`, `ValuesAny`, `CopyAny`, `GetValueType`, `GetIndexType`, etc.) purely so it can satisfy `dataframe.SeriesInterface` / `dataframe.NumericSeriesInterface` — Go generics can't be used across the `Dataframe` package boundary directly, so this boxing-to-`any` layer is how heterogeneous typed columns get stored in a single `DataFrame`.

- **`Dataframe/`** — tabular structure built on top of `series` via interfaces, not concrete types.
  - `seriesInterface.go` defines `SeriesInterface` (what any column must support) and `NumericSeriesInterface` (adds `SumFloat`/`Mean`/`MinFloat`/`MaxFloat`/`Count`/`StdDev`, required for numeric aggregation). `Dataframe` never imports `series` directly — it only depends on these interfaces, which `series.Series`/`series.NumericSeries` satisfy structurally.
  - `dataframe.go` holds `DataFrame` (`columns map[string]SeriesInterface` + `columnOrder []string` + `index`/`indexType`/`nrows`) and its transforms (`Select`, `DropColumn`, `Head`/`Tail`, `FilterColumn`, `ApplyToColumn(s)`, `ResetIndex`, `SetIndexFromColumn`, etc.). Whenever a transform needs to rebuild a column from raw `[]any` (e.g. after filtering rows), it goes through `createSeriesFromAny`, which picks a numeric-capable wrapper (`genericNumericSeries`) when the boxed values are all numeric, or a plain `genericSeries` otherwise — this is what lets a derived/filtered/grouped column still support `GetNumericColumn`.
  - `groupby.go` implements `GroupBy`/`DataFrameGroupBy` and aggregation via `AggFunc` (`SeriesInterface -> any`); `AggregateColumns(map[colName]map[aggSuffix]AggFunc)` is the supported aggregation entry point (result columns are named `col_aggsuffix`). Built-in `AggFunc`s (`AggSum`, `AggMean`, `AggMin`, `AggMax`, `AggStdDev`, `AggCount`, `AggFirst`, `AggLast`) type-assert to `NumericSeriesInterface` where relevant.
  - Because `DataFrame` stores columns by reference (no defensive copy at construction), be careful with code that mutates a `SeriesInterface` in place after handing it to a `DataFrame`.

`main.go` contains three example flows (`example1SumColumn`, `example2AddColumns`, `example3Multipleoperations`) showing the intended usage pattern: build `series`, assemble into a `DataFrame`, retrieve typed numeric columns back out via `GetNumericColumn` + a type assertion to the concrete `*series.NumericSeries[T, R]` when you need type-specific ops (e.g. `Add`/`Multiply`).
