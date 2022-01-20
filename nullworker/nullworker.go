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
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    *zap.Logger
	WG        sync.WaitGroup
}

func (w *NullWorker) Ctx() context.Context {
	if w.ctx == nil {
		w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	}
	return w.ctx
}

func (w *NullWorker) Logger() *zap.Logger {
	return w.logger
}

func (w *NullWorker) Init(l *zap.Logger) error {
	if w.ctx == nil {
		w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	}
	w.logger = l
	return nil
}
func (w *NullWorker) Terminate() error {
	w.ctxCancel()
	w.WG.Wait()
	return nil
}
func (w *NullWorker) Run() error {
	<-w.ctx.Done()
	return nil
}
