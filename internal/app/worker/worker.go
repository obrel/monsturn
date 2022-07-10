package worker

import (
	"context"
	"sync"

	"github.com/obrel/go-lib/pkg/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	tasks     map[string]Task
	numWorker int
	lock      sync.RWMutex
}

func NewWorker(name string) *Worker {
	w := &Worker{
		tasks:     map[string]Task{},
		numWorker: 1,
	}

	return w
}

func (w *Worker) Add(key string, t Task) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.tasks[key] = t
}

func (w *Worker) Get(key string) (Task, error) {
	w.lock.RLock()
	t, ok := w.tasks[key]
	w.lock.RUnlock()
	if !ok {
		return nil, errors.Errorf("Task not found for: %s", key)
	}

	return t, nil
}

func (w *Worker) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	wg := &sync.WaitGroup{}

	for key := range w.tasks {
		wg.Add(1)
		g.Go(func() error {
			task, err := w.Get(key)
			if err != nil {
				log.For("worker", "start").Error(err)
				return err
			}

			err = task.Do(ctx)
			if err != nil {
				log.For("worker", "start").Error(err)
				return err
			}

			wg.Done()
			return nil
		})

		log.For("worker", "start").Infof("Starting worker for %s [%d spawned]", key, w.numWorker)
		wg.Wait()
	}

	return g.Wait()
}
