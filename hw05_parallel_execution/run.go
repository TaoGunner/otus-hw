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
	errCount := atomic.Int64{}

	// Создаём и наполняем пул рабочими
	workerPool := make(chan struct{}, n)
	for range n {
		workerPool <- struct{}{}
	}

	for _, task := range tasks {
		// Ждём свободного рабочего
		<-workerPool

		// Если кол-во ошибок превышает m - прекращаем работу
		if errCount.Load() >= int64(m) {
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Выполняем задание, если ошибка - инкрементируем счетчик
			if err := task(); err != nil {
				errCount.Add(1)
			}

			// Возвращаем освободившегося рабочего в пул
			workerPool <- struct{}{}
		}()
	}

	// Ждём завершения горутин
	wg.Wait()

	// Очищаем и закрываем канал
	for len(workerPool) > 0 {
		<-workerPool
	}
	close(workerPool)

	// Если счетчик ошибок дошел до m - возвращаем ошибку
	if errCount.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
