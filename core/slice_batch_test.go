package core

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallelMap(t *testing.T) {
	tests := []struct {
		name        string
		data        []int
		batchSize   int
		maxWorkers  uint
		mapFunc     func(ctx context.Context, batch []int) ([]string, error)
		expected    []string
		expectError bool
	}{
		{
			name:       "Test Successful Mapping",
			data:       []int{1, 2, 3, 4, 5},
			batchSize:  2,
			maxWorkers: 2,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
				}
				return result, nil
			},
			expected:    []string{"1", "2", "3", "4", "5"},
			expectError: false,
		},
		{
			name:       "Test Empty Data",
			data:       []int{},
			batchSize:  2,
			maxWorkers: 2,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
				}
				return result, nil
			},
			expected:    []string{},
			expectError: false,
		},
		{
			name:       "Test Single Batch",
			data:       []int{1, 2, 3},
			batchSize:  3,
			maxWorkers: 1,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
				}
				return result, nil
			},
			expected:    []string{"1", "2", "3"},
			expectError: false,
		},
		{
			name:       "Test Map Function Error",
			data:       []int{1, 2, 3, 4, 5},
			batchSize:  2,
			maxWorkers: 2,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				if len(batch) == 2 && batch[0] == 3 {
					return nil, fmt.Errorf("error processing batch")
				}
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
				}
				return result, nil
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:       "Test Map with Slow Function",
			data:       []int{1, 2, 3, 4, 5, 6},
			batchSize:  2,
			maxWorkers: 3,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				time.Sleep(100 * time.Millisecond) // Simulate slow operation
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
				}
				return result, nil
			},
			expected:    []string{"1", "2", "3", "4", "5", "6"},
			expectError: false,
		},
		{
			name:       "Test Invalid Batch Size",
			data:       []int{1, 2, 3},
			batchSize:  0,
			maxWorkers: 2,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
				}
				return result, nil
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:       "Test Invalid Max Workers",
			data:       []int{1, 2, 3},
			batchSize:  1,
			maxWorkers: 0,
			mapFunc: func(ctx context.Context, batch []int) ([]string, error) {
				var result []string
				for _, num := range batch {
					result = append(result, strconv.Itoa(num))
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
			result, err := ParallelMap(ctx, tt.data, tt.batchSize, tt.maxWorkers, tt.mapFunc)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
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
			expected:    []int{},
			expectError: false,
		},
		{
			name:       "Test Filtering With Error",
			data:       []int{1, 2, 3, 4, 5},
			batchSize:  2,
			maxWorkers: 2,
			filterFunc: func(ctx context.Context, batch []int) ([]int, error) {
				if len(batch) == 2 && batch[0] == 3 {
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
