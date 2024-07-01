package hw05parallelexecution

import (
	"errors"
	"fmt"
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
	chTask := make(chan Task, len(tasks))
	errCount := int64(0)

	// Запустим требуемое количество горутин
	for range n {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Работаем, пока есть задачи
			for task := range chTask {
				fmt.Println(atomic.LoadInt64(&errCount))
				// Если общий счетчик ошибок дошел до m - прекращаем обработку задач
				if errCount := atomic.LoadInt64(&errCount); errCount >= int64(m) {
					break
				}

				// Если задача завершилась ошибкой - инкрементируем счетчик ошибок
				if err := task(); err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	// Заполняем канал задачами и закрываем его
	for _, task := range tasks {
		chTask <- task
	}
	close(chTask)

	// Ждём завершения горутин
	wg.Wait()

	// Если счетчик ошибок дошел до m - возвращаем ошибку
	if errCount >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
