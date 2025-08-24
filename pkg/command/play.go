package command

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func NewPlay(w io.Writer, mpv, playlist string, loop, dry, window bool) *Play {
	return &Play{
		w:        w,
		mpv:      mpv,
		playlist: playlist,
		loop:     loop,
		dry:      dry,
		window:   window,
	}
}

// Play plays the playlist.
type Play struct {
	w        io.Writer
	mpv      string
	playlist string
	loop     bool
	dry      bool
	window   bool
}

var _ Cmd = &Play{}

func (Play) Name() string { return "play" }

func (c *Play) Run(ctx context.Context) error {
	slog.Info("play", slog.Bool("dry", c.dry))
	if c.dry {
		return c.dryRun()
	}

	arg := []string{
		"--no-video",
		"--playlist=" + c.playlist,
	}
	if c.loop {
		arg = append(arg, "--loop-playlist=true")
	}
	if c.window {
		arg = append(arg, "--force-window=yes")
	} else {
		arg = append(arg, "--force-window=no")
	}
	executable, err := exec.LookPath(c.mpv)
	if err != nil {
		return fmt.Errorf("%w: failed to lookup mpv", err)
	}

	slog.Info("execve", slog.String("cmd", strings.Join(
		append([]string{executable}, arg...),
		" ",
	)))
	// replace the process to interact with mpv
	return syscall.Exec(executable, arg, os.Environ())
}

func (c *Play) dryRun() error {
	f, err := os.Open(c.playlist)
	if err != nil {
		return fmt.Errorf("%w: failed to open playlist", err)
	}
	defer f.Close()

	if _, err := io.Copy(c.w, f); err != nil {
		return fmt.Errorf("%w: failed to read playlist", err)
	}
	return nil
}
