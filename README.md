# psp
Parallel slice processing library

# ParallelMap Function

The `ParallelMap` function provides a mechanism to apply a mapping function concurrently over a slice of data in batches. It processes the data in parallel, using a specified maximum number of workers and batch size. The function can be used to perform concurrent data transformations safely and efficiently, while controlling concurrency and batching to avoid overwhelming resources.

## Function Signature

```go
func ParallelMap[T, R any](
    ctx context.Context,
    data []T,
    batchSize int,
    maxw uint,
    mapFunc func(ctx context.Context, batch []T) ([]R, error),
) ([]R, error)
