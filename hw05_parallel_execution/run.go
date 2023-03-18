package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	mt := sync.Mutex{}
	wg := &sync.WaitGroup{}
	errs := int64(0)

	// Создаем каналы
	taskCh := make(chan Task)
	errCh := make(chan struct{})

	// Добавляем в группу кол-во корутин
	wg.Add(n)

	// Пихаем горутины-обработчики
	for i := 0; i < n; i++ {
		go worker(m, wg, errCh, taskCh, &errs, &mt)
	}

	// Синхрониризуем каналы и пихаем таски в канал
	for _, task := range tasks {
		select {
		case <-errCh:
		case taskCh <- task:
		}
	}

	// Тут не разобрался, без этого горутины не завершаются. Хотя из канала прочитаны все таски и на 71 строчке должны выключаться
	close(taskCh)

	wg.Wait()

	select {
	case <-errCh:
		return ErrErrorsLimitExceeded
	default:
	}

	return nil
}

func worker(
	m int,
	wg *sync.WaitGroup,
	errCh chan struct{},
	taskCh chan Task,
	errs *int64,
	mt *sync.Mutex,
) {
	defer wg.Done()

	for {
		select {
		case <-errCh:
			return

		case task, ok := <-taskCh:
			if !ok {
				return
			}

			if err := task(); err != nil {
				mt.Lock()
				*errs++
				if *errs == int64(m) {
					errCh <- struct{}{}
					close(errCh)
				}
				mt.Unlock()
			}
		}
	}
}
