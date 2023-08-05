package svcf

import (
	"io"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/voi-oss/svc"
	"go.uber.org/zap"
)

type SVC struct {
	*svc.SVC
	workers    map[string]svc.Worker
	flagGroups map[string]interface{}

	// options
	flagHelpOut io.Writer
}

func New(s *svc.SVC, options ...SVCFOption) *SVC {
	svcRes := &SVC{
		SVC:        s,
		workers:    map[string]svc.Worker{},
		flagGroups: map[string]interface{}{},

		flagHelpOut: os.Stdout,
	}

	for _, option := range options {
		option(svcRes)
	}

	return svcRes
}

// Adds a worker which also gets flags parsing for config management through
// github.com/jessevdk/go-flags
func (s *SVC) AddWorker(name string, w svc.Worker) {
	s.workers[name] = w
	s.SVC.AddWorker(name, w)
}

// AddFlagGroup adds a flag group to the service without it requiring to be a worker.
// useful to get consistent flag parsing for multiple modules.
func (s *SVC) AddFlagGroup(name string, fg interface{}) {
	s.flagGroups[name] = fg
}

// Run runs the service until either receiving an interrupt or a worker
// terminates.
func (s *SVC) Run() {
	parser := flags.NewNamedParser(s.SVC.Name, flags.Default)
	for name, w := range s.workers {
		if _, err := parser.AddGroup(name, "", w); err != nil {
			s.Logger().Error("flagparsing", zap.String("modudle_name", name), zap.Error(err))
		}
		// Needed to parse all instances of the same config; otherwise only the last object
		// containing a flag entry will be populated.
		if _, err := parser.Parse(); err != nil {
			s.Logger().Error("flagparsing", zap.String("modudle_name", name), zap.Error(err))
		}
	}

	for name, w := range s.flagGroups {
		if _, err := parser.AddGroup(name, "", w); err != nil {
			s.Logger().Error("flagparsing", zap.String("modudle_name", name), zap.Error(err))
		}
		// Needed to parse all instances of the same config; otherwise only the last object
		// containing a flag entry will be populated.
		if _, err := parser.Parse(); err != nil {
			s.Logger().Error("flagparsing", zap.String("modudle_name", name), zap.Error(err))
		}
	}
	_, err := parser.Parse()
	parser.WriteHelp(s.flagHelpOut)
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
