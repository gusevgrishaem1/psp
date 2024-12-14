package core

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// ParallelMap performs parallel processing on a slice of data using
// batching and limits the number of concurrently running goroutines.
// The function applies the mapFunc to each data batch and returns the final slice
// with the results of all processed batches.
//
// The context (context.Context) is used to manage the execution time
// and cancellation of operations, especially when dealing with long-running
// or parallel tasks.
//
// Parameters:
//   - ctx: the context used to manage cancellation and deadlines.
//   - data: the slice of input data to be processed.
//   - batchSize: the size of each data batch for processing.
//   - maxWorkers: the maximum number of goroutines that can run concurrently.
//   - mapFunc: the function that processes each data batch. It takes
//     the context and a data batch as input and returns a slice of results and an error.
//
// Returns:
// - The resulting slice with processed elements.
// - An error if the processing failed at any stage or if the context was canceled.
//
// Example usage:
// ```go
// ctx := context.Background()
// data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//
//	processFunc := func(ctx context.Context, batch []int) ([]int, error) {
//	    var result []int
//	    for _, v := range batch {
//	        result = append(result, v*2)
//	    }
//	    return result, nil
//	}
//
// result, err := ParallelMap(ctx, data, 3, 3, processFunc)
//
//	if err != nil {
//	    log.Fatalf("Error processing data: %v", err)
//	}
//
// ```
func ParallelMap[T, R any](ctx context.Context, data []T, batchSize int, maxWorkers uint, mapFunc func(ctx context.Context, batch []T) ([]R, error)) ([]R, error) {
	if batchSize <= 0 {
		return nil, errors.New("batch size must be greater than 0")
	}

	if maxWorkers <= 0 {
		return nil, errors.New("maxw must be greater than 0")
	}

	if len(data) == 0 {
		return []R{}, nil
	}

	result := make([]R, len(data))
	sem := make(chan struct{}, maxWorkers)
	defer close(sem)
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < len(data); i += batchSize {
		start := i
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }()

			mappedBatch, err := mapFunc(ctx, data[start:end])
			if err != nil {
				return err
			}

			i := 0
			for j := start; j < end; j++ {
				result[j] = mappedBatch[i]
				i++
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}