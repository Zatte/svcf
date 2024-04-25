package nullworker

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/voi-oss/svc"
)

var _ svc.Worker = (*NullWorker)(nil)

// SNullWorker implements a silent nullworker where all interfaces are implemented
// so no errors are shown during startup.
type SNullWorker struct {
	*NullWorker
}

func NewSNullWorker() *SNullWorker {
	return &SNullWorker{
		NullWorker: New(),
	}
}

func (w *SNullWorker) Healthy() error {
	return nil
}

func (w *SNullWorker) Alive() error {
	return nil
}

func (w *SNullWorker) Gatherer() prometheus.Gatherer {
	return prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) {
		return nil, nil
	})
}
