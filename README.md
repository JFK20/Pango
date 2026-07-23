# Pango

Pango is a Go library that brings pandas-like data manipulation capabilities to Go. The name is a combination of "pandas" and "Go".

## Overview

This project aims to provide familiar data structures and operations for data analysis in Go, inspired by pandas. It leverages Go's generics to offer type-safe data manipulation.
My goal is to port the most important features of pandas to Go, while keeping close to 100% test coverage.

## Current Features

### Series

`Series[T comparable, R comparable]` is a one-dimensional labeled array capable of holding any comparable data type, with a generic label/index type. Built via `NewSeries`/`NewIndexSeries` (default `0..n-1` index). Supports label- and position-based access (`Get`, `At`, `AtIndex`), `Head`/`Tail`, `Append`/`Prepend`, `Copy`, `IsIn`, `ResetIndex`, and sorting by index or by value (`SortByIndex`, `SortByValue`).

### Numeric Series

`NumericSeries[T Numeric, R comparable]` embeds `Series` and adds mathematical operations for numeric data: `Sum`, `Mean`, `Min`/`Max`, `StdDev`, `Quantile`/`Median`, `CoVariance`/`Correlation`, element-wise `Add`/`Subtract`/`Multiply`/`Divide`/`Mod`/`Pow`, `CumSum`, `DropNA`, `ArgMin`/`ArgMax`.

### DataFrame

A tabular structure of named `Series` columns (`Dataframe` package), built via `NewDataFrame`/`NewDataFrameWithIndex`. Supports column management (`AddColumn`, `DropColumn`, `RenameColumn`, `Select`), row access (`Head`, `Tail`, `Range`), row filtering (`FilterColumn`, `FilterIn`), transforms (`ApplyToColumn`, `ApplyToColumns`, `AstypeColumn`), and index manipulation (`ResetIndex`, `SetIndexFromColumn`). Numeric columns are reachable via `GetNumericColumn` without a concrete type assertion.

### GroupBy

Group a DataFrame by one or more columns with `GroupBy`/`GroupByColumns`, then either:
- collapse each group to a single row with `AggregateColumns` (built-in `AggSum`, `AggMean`, `AggMin`, `AggMax`, `AggStdDev`, `AggQuantile`, `AggMedian`, `AggFirst`, `AggLast`, `AggCount`) or `Count`, or
- work with the raw per-group rows directly via `Groups()` (`map[any]*DataFrame`) or the `Range()` iterator.

## Status

This project is in early development. More features and data structures are planned for future releases.
