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
```

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/yourusername/psp" // Replace with your actual package import path
)

func main() {
	ctx := context.Background()
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	batchSize := 3
	maxWorkers := uint(3)

	// Define a filtering function that selects even numbers
	filterFunc := func(ctx context.Context, batch []int) ([]int, error) {
		var result []int
		for _, num := range batch {
			if num%2 == 0 {
				result = append(result, num)
			}
		}
		return result, nil
	}

	// Call ParallelFilter to filter the data in parallel
	result, err := psp.ParallelFilter(ctx, data, batchSize, maxWorkers, filterFunc)
	if err != nil {
		log.Fatalf("Error filtering data: %v", err)
	}

	// Print the results
	fmt.Println("Filtered results:", result)
}
```
