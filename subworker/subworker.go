package subworker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

var _ svc.Worker = (*Orchestrator)(nil)

// Orchestrator is a subclass of Nullworker which allows this worker to take on
// child workers and initialize them much like SVC does but in a slightly
// simpler way.
// The base worker is initialized and run last and shutdown first.
type Orchestrator struct {
	svc.Worker

	ctx       context.Context
	ctxCancel context.CancelFunc

	logger *zap.Logger

	TerminationGracePeriod time.Duration
	TerminationWaitPeriod  time.Duration

	workers            map[string]svc.Worker
	workersAdded       []string
	workersInitialized []string

	terminateErrors chan error
}

func NewOrchestrator(baseWorker svc.Worker) *Orchestrator {
	ctx, ctxCancel := context.WithCancel(context.Background())
	res := &Orchestrator{
		Worker:    baseWorker,
		ctx:       ctx,
		ctxCancel: ctxCancel,

		logger: zap.L(),

		workers:            map[string]svc.Worker{},
		workersAdded:       []string{},
		workersInitialized: []string{},

		terminateErrors: make(chan error),
	}
	return res
}

func (s *Orchestrator) AddSubWorker(name string, w svc.Worker) error {
	if len(s.workersInitialized) > 0 {
		return fmt.Errorf("sub worker %s added after initialization", name)
	}
	if _, exists := s.workers[name]; exists {
		return fmt.Errorf("sub worker with name %s added twice", name)
	}
	// Track workers as ordered set to initialize them in order.
	s.workersAdded = append(s.workersAdded, name)
	s.workers[name] = w

	return nil
}

func (s *Orchestrator) Init(l *zap.Logger) error {
	// Initializing workers in added order.
	s.logger = l
	s.AddSubWorker("_", s.Worker) // The baseworker should be initialized last.
	for _, name := range s.workersAdded {
		l := l
		if name != "_" {
			l.Debug("Initializing sub worker", zap.String("sub_worker", name))
			l = l.Named(name)
		}
		w := s.workers[name]

		if err := w.Init(l); err != nil {
			l.Error("Could not initialize sub worker", zap.String("sub_worker", name), zap.Error(err))

			return err
		}
		s.workersInitialized = append(s.workersInitialized, name)
	}

	return nil
}

func (s *Orchestrator) recoverWait(name string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Done()
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			s.logger.Error("recover panic", zap.String("sub_worker", name),
				zap.Error(err), zap.Stack("stack"))
			errors <- err
		} else {
			errors <- fmt.Errorf("%v", r)
		}
	}
}

func (s *Orchestrator) Run() error {
	defer s.terminate()

	// Initializing workers in added order.
	errs := make(chan error)
	wg := sync.WaitGroup{}
	for name, w := range s.workers {
		wg.Add(1)
		go func(name string, w svc.Worker) {
			defer s.recoverWait(name, &wg, errs)
			if err := w.Run(); err != nil {
				err = fmt.Errorf("sub worker %s exited: %w", name, err)
				errs <- err
			}
		}(name, w)
	}

	select {
	case err := <-errs:
		s.logger.Fatal("Worker Init/Run failure", zap.Error(err))
		return err
	case <-s.ctx.Done():
		s.logger.Warn("Worker termination requested")
	case <-waitGroupToChan(&wg):
		s.logger.Info("All sub workers have finished")
	}

	return nil
}

func (s *Orchestrator) Terminate() error {
	s.ctxCancel()
	return <-s.terminateErrors
}

func (s *Orchestrator) terminate() {
	s.logger.Info("Terminating sub workers")

	defer close(s.terminateErrors)

	multiErrors := []error{}

	// terminate only initialized workers
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, name := range s.workersInitialized {
			defer func(name string) {
				w := s.workers[name]
				if err := w.Terminate(); err != nil {
					multiErrors = append(multiErrors, err)
					if name == "_" {
						if name != "_" {
							s.logger.Error("Worker terminated with error",
								zap.Error(err))
						}
					} else {
						s.logger.Error("Sub worker terminated with error",
							zap.String("sub_worker", name),
							zap.Error(err))
					}
				} else {
					if name != "_" {
						s.logger.Info("Sub worker terminated", zap.String("sub_worker", name))
					}
				}
			}(name)
		}
	}()
	wg.Wait()
	s.logger.Info("All sub workers terminated, Worker shutdown complete")

	if len(multiErrors) == 0 {
		return
	}

	s.terminateErrors <- errors.Join(multiErrors...)
}

func waitGroupToChan(wg *sync.WaitGroup) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	return c
}
