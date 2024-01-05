package sentryworker

import (
	"context"

	"github.com/voi-oss/svc"
	"github.com/zatte/svcf/nullworker"
	"go.uber.org/zap"
)

// SentryWorker allows for limited introspecting of the worker life-cycle.
// this is  mostly useful in testing where you need to know if workers have
// been  initialized/started/terminated.

type SentryWorker struct {
	svc.Worker
	initDone chan struct{}
	initErr  error

	runCalled       chan struct{}
	terminateCalled chan struct{}

	runCompleted chan struct{}
	runErr       error

	terminateCompleted chan struct{}
	terminateErr       error
}

// Creates a new sentry. Can wrap an existing worker, if nil is provided a
// nullworker is created.
func New(maybeWorker svc.Worker) *SentryWorker {
	if maybeWorker == nil {
		maybeWorker = nullworker.NewNullWorker()
	}
	return &SentryWorker{
		Worker:   maybeWorker,
		initDone: make(chan struct{}),

		runCalled:    make(chan struct{}),
		runCompleted: make(chan struct{}),

		terminateCalled:    make(chan struct{}),
		terminateCompleted: make(chan struct{}),
	}
}

func (w *SentryWorker) Init(l *zap.Logger) error {
	defer close(w.initDone)

	w.initErr = w.Worker.Init(l) //nolint:wrapcheck
	return w.initErr
}

func (w *SentryWorker) Run() error {
	close(w.runCalled)
	defer close(w.runCompleted)

	w.runErr = w.Worker.Run() //nolint:wrapcheck
	return w.runErr
}

func (w *SentryWorker) Terminate() error {
	close(w.terminateCalled)
	defer close(w.terminateCompleted)

	w.terminateErr = w.Worker.Terminate() //nolint:wrapcheck
	return w.terminateErr
}

// InitDone checks if Init() has been called
func (w *SentryWorker) InitDone() bool {
	select {
	case <-w.initDone:
		return true
	default:
		return false
	}
}

// RunIsCalled checks if Run() has been called (but might not have) completed
func (w *SentryWorker) RunIsCalled() bool {
	select {
	case <-w.runCalled:
		return true
	default:
		return false
	}
}

// RunIsCompleted checks if Run() has completed
func (w *SentryWorker) RunIsCompleted() bool {
	select {
	case <-w.runCompleted:
		return true
	default:
		return false
	}
}

// TerminatedIsCalled checks if Terminate() has been called (but might not have) completed
func (w *SentryWorker) TerminatedIsCalled() bool {
	select {
	case <-w.terminateCalled:
		return true
	default:
		return false
	}
}

// TerminatedIsCompleted checks if Terminate() has completed
func (w *SentryWorker) TerminatedIsCompleted() bool {
	select {
	case <-w.terminateCompleted:
		return true
	default:
		return false
	}
}

// WaitForInitDone blocks until Init() is done or the context is canceled.
// The initialization error is returned (may be nil) or context.Cause(ctx) is returned
func (w *SentryWorker) WaitForInitDone(ctx context.Context) error {
	select {
	case <-w.initDone:
		return w.initErr
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

// WaitForRunCalled blocks until Run() is started or the context is canceled.
// If the context is canceled, the error is returned; otherwise nil
func (w *SentryWorker) WaitForRunCalled(ctx context.Context) error {
	select {
	case <-w.runCalled:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx) //nolint:wrapcheck
	}
}

// WaitForRunCompleted blocks until Run() has completed
// The Run() error (may be nil) or context.Cause(ctx) is returned
func (w *SentryWorker) WaitForRunCompleted(ctx context.Context) error {
	select {
	case <-w.runCompleted:
		return w.runErr
	case <-ctx.Done():
		return context.Cause(ctx) //nolint:wrapcheck
	}
}

// WaitForTerminateCalled blocks until Terminate() is started or the context is canceled.
// If the context is canceled, the error is returned; otherwise nil
func (w *SentryWorker) WaitForTerminateCalled(ctx context.Context) error {
	select {
	case <-w.terminateCalled:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx) //nolint:wrapcheck
	}
}

// WaitForTerminateCompleted blocks until Run() is started or the context is canceled.
// The Terminate() error (may be nil) or context.Cause(ctx) is returned
func (w *SentryWorker) WaitForTerminateCompleted(ctx context.Context) error {
	select {
	case <-w.terminateCompleted:
		return w.terminateErr
	case <-ctx.Done():
		return context.Cause(ctx) //nolint:wrapcheck
	}
}
