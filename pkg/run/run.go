package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/berquerant/local-jukebox/pkg/command"
	"github.com/berquerant/local-jukebox/pkg/config"
	"github.com/berquerant/local-jukebox/pkg/iox"
)

func Main(c *config.Config) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGPIPE,
	)
	defer stop()
	r := &runner{
		c: c,
	}
	return r.run(ctx)
}

type runner struct {
	c          *config.Config
	existIndex bool
}

func (r *runner) run(ctx context.Context) error {
	if err := r.c.Validate(); err != nil {
		return err
	}

	if r.c.PlayOnly {
		dest, err := iox.NewTmpFile(os.Stdin, "playlist")
		if err != nil {
			return err
		}
		return r.runPlay(ctx, dest.Name())
	}

	r.existIndex = iox.ExistFile(r.c.IndexFile())

	if err := r.runIndex(ctx); err != nil {
		return err
	}
	if err := r.runNormalize(ctx); err != nil {
		return err
	}
	dest, err := r.runQuery(ctx)
	if err != nil {
		return err
	}
	if err := r.runPlay(ctx, dest.Name()); err != nil {
		return err
	}
	return nil
}

func (r *runner) runIndex(ctx context.Context) error {
	if !r.c.Reload && r.existIndex {
		return nil
	}

	mc, err := r.c.MetafindConfig()
	if err != nil {
		return err
	}
	c := command.NewIndex(
		r.c.MetafindCmd,
		mc,
		r.c.MusicRoot,
		r.c.IndexFile(),
	)
	return command.Run(ctx, c)
}

func (r *runner) runNormalize(ctx context.Context) error {
	if !r.c.Reload && !r.c.Normalize && r.existIndex {
		return nil
	}

	c := command.NewNormalize(
		r.c.NormalizedIndexFile(),
		r.c.IndexFile(),
	)
	return command.Run(ctx, c)
}

func (r *runner) runQuery(ctx context.Context) (*iox.TmpFile, error) {
	var emptyBuf bytes.Buffer
	dest, err := iox.NewTmpFile(&emptyBuf, "playlist")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create playlist", err)
	}

	var query io.Reader = os.Stdin
	if r.c.Query != "stdin" {
		f, err := os.Open(r.c.Query)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to open query", err)
		}
		defer f.Close()
		query = f
	}

	c := command.NewQuery(
		r.c.MetafindCmd,
		r.c.GrepCmd,
		query,
		r.c.NormalizedIndexFile(),
		r.c.ListArgs,
		r.c.GrepArgs,
		r.c.Lines,
		r.c.Shuffle,
		dest.Name(),
	)
	if err := command.Run(ctx, c); err != nil {
		return nil, fmt.Errorf("%w: failed to execute query", err)
	}
	return dest, nil
}

func (r *runner) runPlay(ctx context.Context, playlist string) error {
	c := command.NewPlay(
		r.c.Writer,
		r.c.MpvCmd,
		playlist,
		r.c.Loop,
		r.c.Dry,
		r.c.Window,
	)
	return command.Run(ctx, c)
}
