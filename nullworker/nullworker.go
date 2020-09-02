package nullworker

import (
	"context"

	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

var _ svc.Worker = (*NullWorker)(nil)

// NullWorker implements a minimal worker controled with contexts.
type NullWorker struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (w *NullWorker) Ctx() context.Context {
	return w.ctx
}

func (w *NullWorker) Init(*zap.Logger) error {
	w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	return nil
}
func (w *NullWorker) Terminate() error {
	<-w.ctx.Done()
	return nil
}
func (w *NullWorker) Run() error {
	w.ctxCancel()
	return nil
}
