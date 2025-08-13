package config

import (
	"io"
	"log/slog"
)

func (c *Config) SetupLogger(w io.Writer) {
	level := slog.LevelInfo
	if c.Debug {
		level = slog.LevelDebug
	}
	if c.Quiet {
		level = slog.LevelError
	}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}
