package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded    = errors.New("errors limit exceeded")
	ErrInvalidGoroutinesCount = errors.New("invalid goroutines count")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Если количество горутин <= 0, значит вернуть ErrInvalidGoroutinesCount
	if n <= 0 {
		return ErrInvalidGoroutinesCount
	}
	// "максимум 0 ошибок", значит функция всегда будет возвращать ErrErrorsLimitExceeded
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	chTask := make(chan Task)
	errCount := atomic.Int64{}

	// Запустим n горутин
	for range n {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Работаем, пока есть задачи
			for task := range chTask {
				// Если задача завершилась ошибкой - инкрементируем счетчик ошибок
				if err := task(); err != nil {
					errCount.Add(1)
				}
			}
		}()
	}

	// Заполняем канал задачами
	for _, task := range tasks {
		// Если кол-во ошибок превышает m - прекращаем подавать задания
		if errCount.Load() >= int64(m) {
			break
		}
		chTask <- task
	}
	close(chTask)

	// Ждём завершения горутин
	wg.Wait()

	// Если счетчик ошибок дошел до m - возвращаем ошибку
	if errCount.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
