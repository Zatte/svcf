package svcf

import "io"

type SVCFOption func(s *SVC)

// stops printing flag help on startup.
func WithNoFlagHelp() SVCFOption {
	return WithFlagHelpOut(io.Discard)
}

// direct flag based help messages to this writer.
func WithFlagHelpOut(w io.Writer) SVCFOption {
	return func(s *SVC) {
		s.flagHelpOut = w
	}
}
