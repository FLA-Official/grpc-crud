package utils

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type PrettyHandler struct{}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {

	fmt.Println("[LOG]")
	fmt.Printf("\"time\": \"%s\",\n", r.Time.Format(time.RFC3339))

	level := r.Level.String()
	fmt.Printf("\"level\": \"%s\",\n", level)

	fmt.Printf("\"msg\": \"%s\",\n", r.Message)

	r.Attrs(func(a slog.Attr) bool {
		fmt.Printf("\"%s\": \"%v\",\n", a.Key, a.Value)
		return true
	})

	fmt.Println()

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return h
}

func NewLogger() *slog.Logger {
	return slog.New(&PrettyHandler{})
}
