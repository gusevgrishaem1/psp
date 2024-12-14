# PSP: Parallel Slice Processing Library

The **PSP** (Parallel Slice Processing) library provides a set of utilities for processing slices of data in parallel, enabling efficient concurrent operations while controlling concurrency and resource utilization.

## ParallelMap Function

The `ParallelMap` function allows you to apply a transformation function (mapping function) concurrently over a slice of data in batches. This function splits the data into smaller chunks, processes them in parallel, and then combines the results. By controlling the number of workers (`maxw`) and the batch size (`batchSize`), you can fine-tune the function to suit your performance and resource requirements.

### Function Signature

```go
func ParallelMap[T, R any](
    ctx context.Context,
    data []T,
    batchSize int,
    maxw uint,
    mapFunc func(ctx context.Context, batch []T) ([]R, error),
) ([]R, error)
```

## ParallelFilter Function

The `ParallelFilter` function allows you to filter data concurrently in batches. It applies a filtering function over the data, retaining only the elements that meet the filter condition. This function enables efficient parallel processing while maintaining control over concurrency and batch sizes.

### Function Signature

```go
func ParallelFilter[T any](
    ctx context.Context,
    data []T,
    batchSize int,
    maxw uint,
    filterFunc func(ctx context.Context, batch []T) ([]T, error),
) ([]T, error)
