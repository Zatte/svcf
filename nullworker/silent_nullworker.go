package nullworker

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

var _ svc.Worker = (*NullWorker)(nil)

// SNullWorker implements a silent nullworker where all interfaces are implemented
// so no errors are shown during startup.
type SNullWorker struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	logger *zap.Logger
}

func (w *SNullWorker) Ctx() context.Context {
	return w.ctx
}

func (w *SNullWorker) Logger() *zap.Logger {
	return w.logger
}

func (w *SNullWorker) Init(l *zap.Logger) error {
	w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	w.logger = l
	return nil
}
func (w *SNullWorker) Terminate() error {
	w.ctxCancel()
	return nil
}
func (w *SNullWorker) Run() error {
	<-w.ctx.Done()
	return nil
}

func (w *SNullWorker) Healthy() error {
	return nil
}

func (w *SNullWorker) Gatherer() prometheus.Gatherer {
	return prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) {
		return nil, nil
	})
}
