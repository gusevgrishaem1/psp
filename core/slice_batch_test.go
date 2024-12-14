package core

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock map function that simply returns each element as a string
func mockMapFunc(ctx context.Context, batch []int) ([]string, error) {
	var result []string
	for _, num := range batch {
		result = append(result, strconv.Itoa(num))
	}
	return result, nil
}

func TestParallelMap_Success(t *testing.T) {
	ctx := context.Background()

	// Test input
	data := []int{1, 2, 3, 4, 5}
	batchSize := 2
	maxw := uint(2) // Max 2 workers concurrently

	// Define expected output
	expected := []string{"1", "2", "3", "4", "5"}

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParallelMap_EmptyData(t *testing.T) {
	ctx := context.Background()

	// Test empty input
	data := []int{}
	batchSize := 2
	maxw := uint(2) // Max 2 workers concurrently

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

	// Assertions
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestParallelMap_SingleBatch(t *testing.T) {
	ctx := context.Background()

	// Test input that fits in a single batch
	data := []int{1, 2, 3}
	batchSize := 3
	maxw := uint(1) // Max 1 worker concurrently

	// Expected output
	expected := []string{"1", "2", "3"}

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParallelMap_Errors(t *testing.T) {
	ctx := context.Background()

	// Mock map function that returns an error
	mockErrMapFunc := func(ctx context.Context, batch []int) ([]string, error) {
		return nil, errors.New("batch processing error")
	}

	// Test input that will trigger an error in the map function
	data := []int{1, 2, 3, 4, 5}
	batchSize := 2
	maxw := uint(2)

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockErrMapFunc)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "batch processing error", err.Error())
}

func TestParallelMap_Concurrency(t *testing.T) {
	ctx := context.Background()

	// Test input data
	data := []int{1, 2, 3, 4, 5, 6}
	batchSize := 2
	maxw := uint(3) // Max 3 workers concurrently

	// Simulating slow map function to check if concurrency works
	mockSlowMapFunc := func(ctx context.Context, batch []int) ([]string, error) {
		// Simulate processing delay
		time.Sleep(100 * time.Millisecond)
		var result []string
		for _, num := range batch {
			result = append(result, string(rune(num)))
		}
		return result, nil
	}

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockSlowMapFunc)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, len(data)) // Length should match input length
}

func TestParallelMap_InvalidBatchSize(t *testing.T) {
	ctx := context.Background()

	// Invalid batch size
	data := []int{1, 2, 3}
	batchSize := 0
	maxw := uint(2)

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "batch size must be greater than 0", err.Error())
}

func TestParallelMap_InvalidMaxWorkers(t *testing.T) {
	ctx := context.Background()

	// Invalid max workers
	data := []int{1, 2, 3}
	batchSize := 1
	maxw := uint(0)

	// Execute ParallelMap
	result, err := ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "maxw must be greater than 0", err.Error())
}

func TestParallelFilter(t *testing.T) {
	tests := []struct {
		name        string
		data        []int
		batchSize   int
		maxWorkers  uint
		filterFunc  func(ctx context.Context, batch []int) ([]int, error)
		expected    []int
		expectError bool
	}{
		{
			name:       "Test Filtering Even Numbers",
			data:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			batchSize:  3,
			maxWorkers: 3,
			filterFunc: func(ctx context.Context, batch []int) ([]int, error) {
				var result []int
				for _, v := range batch {
					if v%2 == 0 {
						result = append(result, v)
					}
				}
				return result, nil
			},
			expected:    []int{2, 4, 6, 8, 10},
			expectError: false,
		},
		{
			name:       "Test Empty Data",
			data:       []int{},
			batchSize:  3,
			maxWorkers: 3,
			filterFunc: func(ctx context.Context, batch []int) ([]int, error) {
				var result []int
				for _, v := range batch {
					if v%2 == 0 {
						result = append(result, v)
					}
				}
				return result, nil
			},
			expected:    []int{},
			expectError: false,
		},
		{
			name:       "Test Data With No Even Numbers",
			data:       []int{1, 3, 5, 7, 9},
			batchSize:  2,
			maxWorkers: 2,
			filterFunc: func(ctx context.Context, batch []int) ([]int, error) {
				var result []int
				for _, v := range batch {
					if v%2 == 0 {
						result = append(result, v)
					}
				}
				return result, nil
			},
			expected:    []int{}, // No even numbers
			expectError: false,
		},
		{
			name:       "Test Filtering With Error",
			data:       []int{1, 2, 3, 4, 5},
			batchSize:  2,
			maxWorkers: 2,
			filterFunc: func(ctx context.Context, batch []int) ([]int, error) {
				if len(batch) == 2 && batch[0] == 3 { // Simulate error on 3
					return nil, fmt.Errorf("error processing batch")
				}
				var result []int
				for _, v := range batch {
					if v%2 == 0 {
						result = append(result, v)
					}
				}
				return result, nil
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := ParallelFilter(ctx, tt.data, tt.batchSize, tt.maxWorkers, tt.filterFunc)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
