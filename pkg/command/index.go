package command

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/berquerant/local-jukebox/pkg/iox"
)

func NewIndex(metafind, metafindConfig, musicRoot, dest string) *Index {
	return &Index{
		metafind:       metafind,
		metafindConfig: metafindConfig,
		musicRoot:      musicRoot,
		dest:           dest,
	}
}

// Index creates the metafind index file.
type Index struct {
	metafind       string
	metafindConfig string
	musicRoot      string
	dest           string
}

var _ Cmd = &Index{}

func (Index) Name() string { return "index" }

func (c *Index) Run(ctx context.Context) error {
	tmpf, err := iox.NewTmpFile(bytes.NewBufferString(c.metafindConfig), "mf.yml")
	if err != nil {
		return fmt.Errorf("%w: failed to write metafind config", err)
	}
	defer tmpf.Close()

	dest, err := os.Create(c.dest)
	if err != nil {
		return fmt.Errorf("%w: failed create dest file", err)
	}
	defer dest.Close()

	cmd := exec.CommandContext(ctx, c.metafind, "--config", tmpf.Name(), "-v")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "MUSIC_ROOT="+c.musicRoot) // ensure MUSIC_ROOT
	cmd.Stdout = dest
	cmd.Stderr = os.Stderr
	slog.Info("index", slog.String("root", c.musicRoot), slog.String("dest", c.dest))
	logCmd(cmd)
	return cmd.Run()
}
