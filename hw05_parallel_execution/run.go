package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func executeTask(taskChan chan Task,
	wg *sync.WaitGroup,
	stopChan chan struct{},
	m int,
	errCount *int64,
	once *sync.Once,
) {
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			return
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			if err := task(); err != nil && m > 0 {
				if atomic.AddInt64(errCount, 1) >= int64(m) {
					once.Do(func() { close(stopChan) })
					return
				}
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	if n <= 0 {
		return fmt.Errorf("invalid workers count")
	}

	taskChan := make(chan Task, min(len(tasks), n))
	stopChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(n)

	var errCount int64

	// worker gorutines
	once := sync.Once{}
	for i := 0; i < n; i++ {
		go executeTask(taskChan, &wg, stopChan, m, &errCount, &once)
	}

	// fill task channel with tasks
	for _, taskFunc := range tasks {
		shouldEnd := false
		select {
		case <-stopChan:
			shouldEnd = true
		case taskChan <- taskFunc:
			continue
		}
		if shouldEnd {
			break
		}
	}

	close(taskChan)

	wg.Wait()

	if m > 0 && atomic.LoadInt64(&errCount) >= int64(m) {
		once.Do(func() { close(stopChan) })
		return ErrErrorsLimitExceeded
	}
	close(stopChan)
	return nil
}
