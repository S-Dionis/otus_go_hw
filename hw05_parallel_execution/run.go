package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		return errors.New("n must be > 0")
	}

	if m <= 0 {
		m = len(tasks) + 1
	}

	tasksChan := make(chan Task)
	var errorCount int32

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				if task() != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorCount) >= int32(m) {
			break
		}
		tasksChan <- task
	}

	close(tasksChan)
	wg.Wait()

	if errorCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
