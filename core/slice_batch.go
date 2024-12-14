package core

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// ParallelMap выполняет параллельную обработку среза данных с использованием
// пакетов (batching) и ограничением на количество одновременно работающих горутин.
// Функция применяет mapFunc к каждому пакету данных и возвращает итоговый срез
// с результатами всех обработанных пакетов.
//
// Контекст (context.Context) используется для управления временем выполнения
// и отмены операций, особенно при длительных или параллельных задачах.
//
// Параметры:
//   - ctx: контекст, используемый для управления отменой и дедлайнами.
//   - data: срез входных данных, который будет обрабатываться.
//   - batchSize: размер каждого пакета данных для обработки.
//   - maxw: максимальное количество горутин, которые могут работать одновременно.
//   - mapFunc: функция, которая обрабатывает каждый пакет данных. Она принимает
//     контекст и пакет данных и возвращает срез результатов и ошибку.
//
// Возвращает:
// - Результирующий срез с обработанными элементами.
// - Ошибку, если обработка не удалась на любом этапе или если контекст был отменен.
//
// Пример использования:
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
