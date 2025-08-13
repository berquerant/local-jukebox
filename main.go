package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/berquerant/local-jukebox/pkg/config"
	"github.com/berquerant/local-jukebox/pkg/run"
	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

const usage = `jukebox - play music files by querying your local library

# Usage

  jukebox [flags] [mf -i additional args...] [--] [grep args...]

# Examples

Play music files with the query file.

  jukebox -r /root/dir/of/music -x query.txt

Display music files.

  MUSIC_ROOT=/root/dir/of/music jukebox --dry < query.txt

Reload the index and display music files.

  jukebox -r /root/dir/of/music -x query.txt --dry --reload

Limit the music file count to 3.

  jukebox -r /root/dir/of/music -x query.txt --dry -n 3

Grep the music files.

  jukebox -r /root/dir/of/music -x query.txt --dry -- keyword

Loop.

  jukebox -r /root/dir/of/music -x query.txt --loop

# Prerequisites

- mf https://github.com/berquerant/metafind
- mpv https://github.com/mpv-player/mpv
- grep
- jq https://github.com/jqlang/jq
- ffprobe https://ffmpeg.org/ffprobe.html

# Flags`

func main() {
	fs := pflag.NewFlagSet("main", pflag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println(usage)
		fs.PrintDefaults()
	}

	c, err := structconfig.NewConfigWithMerge(
		structconfig.New[config.Config](),
		structconfig.NewMerger[config.Config](),
		fs,
	)
	if errors.Is(err, pflag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		failed()
	}

	c.SetupLogger(os.Stderr)
	c.Writer = os.Stdout

	args := fs.Args() // jukebox POSITIONAL_ARGS...
	dashPos := fs.ArgsLenAtDash()
	if dashPos > 1 { // jukebox ARG... -- [ARG...]
		c.ListArgs = args[1:dashPos]
	}
	if dashPos > 0 { // jukebox [ARG...] -- ARG...
		c.GrepArgs = args[dashPos:]
	}
	if dashPos < 0 { // jukebox [ARG...] (-- is not found)
		c.ListArgs = args[1:]
	}

	cj, _ := json.Marshal(c)
	slog.Debug("config", slog.String("json", string(cj)))
	if err := run.Main(c); err != nil {
		slog.Error("exit", slog.Any("err", err))
		failed()
	}
}

func failed() {
	os.Exit(1)
}
