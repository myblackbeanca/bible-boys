package concurrent

import "sync"

func RunTasks(tasks []func() error, maxConcurrent int) []error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)
	errs := make([]error, len(tasks))
	for i, task := range tasks {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, task func() error) {
			defer wg.Done()
			errs[i] = task()
			<-sem
		}(i, task)
	}
	wg.Wait()
	return errs
}