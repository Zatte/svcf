package nullworker

import (
	"context"

	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

var _ svc.Worker = (*NullWorker)(nil)

// NullWorker implements a minimal worker controled with contexts.
type NullWorker struct {
	Ctx       context.Context
	CtxCancel context.CancelFunc
}

func (w *NullWorker) Init(*zap.Logger) error {
	w.Ctx, w.CtxCancel = context.WithCancel(context.Background())
	return nil
}
func (w *NullWorker) Terminate() error {
	<-w.Ctx.Done()
	return nil
}
func (w *NullWorker) Run() error {
	w.CtxCancel()
	return nil
}
