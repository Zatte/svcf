package svcf

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/voi-oss/svc"
)

type SVC struct {
	svc.SVC
	workers map[string]svc.Worker
}

func New(s svc.SVC) *SVC {
	return &SVC{
		SVC:     s,
		workers: map[string]svc.Worker{},
	}
}

func (s *SVC) AddWorker(name string, w svc.Worker) {
	s.workers[name] = w
	s.SVC.AddWorker(name, w)
}

// Run runs the service until either receiving an interrupt or a worker
// terminates.
func (s *SVC) Run() {
	s.SVC.Logger().Info("Paring flags")

	parser := flags.NewNamedParser(s.SVC.Name, flags.Default)
	for name, w := range s.workers {
		parser.AddGroup(name, "", w)
	}
	_, err := parser.Parse()
	parser.WriteHelp(os.Stdout)
	if err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		s.SVC.Logger().Error(err.Error())
		os.Exit(code)
	}

	s.SVC.Run()
}
