// Package slogger wraps log/slog for more flexibility.
//
// TODO: possibly ditch this because it's a LOT of complexity for just the
// removal of extra attrs maybe we don't even want.
package slogger

// The goal here is to be able to add ReplaceAttr logic to an existing
// logger by cloning it similarly to .With, but .WithReplaceAttr.

import (
	"io"
	"log/slog"
)

type HandlerMaker interface {
	New(w io.Writer, opts *slog.HandlerOptions) slog.Handler
}

type Slogger struct {
	*slog.Logger
	HandlerMaker HandlerMaker
	// Keep here or better to keep elsewhere?
	Writer io.Writer
	Opts   *slog.HandlerOptions
}

// WithReplaceAttr returns a new Slogger build with r, but with the same
// Writer and other Opts as s.
//
// TODO: With(attr ...any) -- needs to keep track of them so we can reproduce.
func (s *Slogger) WithReplaceAttr(r func(groups []string, a slog.Attr) slog.Attr) *Slogger {
	opts := &slog.HandlerOptions{
		AddSource: s.Opts.AddSource,
		Level:     s.Opts.Level,
	}

	return NewSlogger(
		s.HandlerMaker,
		s.Writer,
		opts,
	)

}

func NewSlogger(f HandlerMaker, w io.Writer, opts *slog.HandlerOptions) *Slogger {
	return &Slogger{
		slog.New(f.New(w, opts)),
		f,
		w,
		opts,
	}
}

type JsonHandlerMaker struct {
}

func (JsonHandlerMaker) New(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, opts)
}

type TextHandlerMaker struct {
}

func (TextHandlerMaker) New(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, opts)
}

func demo(w io.Writer, json bool) {

	var maker HandlerMaker
	if json {
		maker = JsonHandlerMaker{}
	} else {
		maker = TextHandlerMaker{}
	}
	logger := NewSlogger(maker, w, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger.Info("orig here at info", "this", "that")
	logger.Debug("orig here at debug", "this", "that")

	// Now clone it with a ReplaceAttr.
	replacer := func(groups []string, a slog.Attr) slog.Attr {
		// Remove foo, for instance!
		if a.Key == "foo" {
			return slog.Attr{}
		}
		return a
	}
	panic(replacer)

}
