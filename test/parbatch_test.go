package test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/gusevgrishaem1/psp/parbatch"
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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockErrMapFunc)

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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockSlowMapFunc)

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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

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
	result, err := parbatch.ParallelMap(ctx, data, batchSize, maxw, mockMapFunc)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "maxw must be greater than 0", err.Error())
}
