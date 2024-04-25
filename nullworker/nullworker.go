package nullworker

import (
	"context"
	"sync"

	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

var _ svc.Worker = (*NullWorker)(nil)

// NullWorker implements a minimal worker controled with contexts.
type NullWorker struct {
	mu        sync.Mutex
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    *zap.Logger
	WG        sync.WaitGroup
}

func New() *NullWorker {
	ctx, ctxCancel := context.WithCancel(context.Background())
	res := &NullWorker{
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
	return res
}

// NewNullWorker is Deprecated: use New instead.
func NewNullWorker() *NullWorker {
	return New()
}

func (w *NullWorker) Ctx() context.Context {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.ctx == nil {
		if w.ctx == nil {
			w.ctx, w.ctxCancel = context.WithCancel(context.Background())
		}
	}
	return w.ctx
}

func (w *NullWorker) Logger() *zap.Logger {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.logger == nil {
		if w.logger == nil {
			l, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			return l
		}
	}
	return w.logger
}

func (w *NullWorker) Init(l *zap.Logger) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.ctx == nil {
		w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	}
	w.logger = l
	return nil
}

func (w *NullWorker) Terminate() error {
	if w.ctxCancel != nil {
		w.ctxCancel()
	}
	w.WG.Wait()
	return nil
}

func (w *NullWorker) Run() error {
	<-w.Ctx().Done()
	return nil
}
