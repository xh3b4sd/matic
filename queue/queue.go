package queue

import (
	"sync"
)

type Queue struct {
	Ctx interface{}
}

type Task func(ctx interface{}) error

func NewQueue(ctx interface{}) *Queue {
	return &Queue{
		Ctx: ctx,
	}
}

func (q *Queue) RunTasks(tasks ...Task) error {
	if len(tasks) == 0 {
		return Mask(ErrNoTasks)
	}

	if err := inSeries(q.Ctx, tasks...); err != nil {
		return Mask(err)
	}

	return nil
}

func (q *Queue) InSeries(tasks ...Task) Task {
	return func(ctx interface{}) error {
		if err := inSeries(ctx, tasks...); err != nil {
			return Mask(err)
		}

		return nil
	}
}

func (q *Queue) InParallel(tasks ...Task) Task {
	return func(ctx interface{}) error {
		// Create error for current tasks. If there occurs one error, all remaining
		// tasks will be canceled.
		var err error

		// Create waitgroup to keep track of parallel tasks by registering the count
		// of them.
		var wg sync.WaitGroup
		wg.Add(len(tasks))

		// Start a goroutine for each task to run them in parallel.
		for i, t := range tasks {
			go func(t Task, ctx interface{}) {
				if e := t(ctx); e != nil {
					err = e

					// Cancel all remaining go routines.
					diff := len(tasks) - i + 1
					for j := 0; j < diff; j++ {
						wg.Done()
					}
				}

				wg.Done()
			}(t, ctx)
		}

		// Wait until the waitgroup count is 0.
		wg.Wait()

		if err != nil {
			return Mask(err)
		}

		return nil
	}
}

func inSeries(ctx interface{}, tasks ...Task) error {
	for _, t := range tasks {
		if err := t(ctx); err != nil {
			return Mask(err)
		}
	}

	return nil
}
