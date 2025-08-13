package command

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

type Cmd interface {
	Name() string
	Run(ctx context.Context) error
}

func Run(ctx context.Context, cmd Cmd) error {
	now := time.Now()
	slog.Debug("start", slog.String("cmd", cmd.Name()))
	defer func() {
		slog.Debug("end", slog.String("cmd", cmd.Name()), slog.String("duration", time.Since(now).String()))
	}()

	if err := cmd.Run(ctx); err != nil {
		return fmt.Errorf("%w: %s", err, cmd.Name())
	}
	return nil
}

func logCmd(cmd *exec.Cmd) {
	slog.Debug("exec", slog.String("cmd", strings.Join(cmd.Args, " ")))
}

func isExitWith(err error, exitCode int) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == exitCode
	}
	return false
}
