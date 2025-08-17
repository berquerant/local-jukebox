package config

import (
	"io"
	"path/filepath"
)

type Config struct {
	//
	// CLI options, envs
	//
	Query     string `name:"query" short:"x" default:"stdin" usage:"music query"`
	Dry       bool   `name:"dry" short:"l" usage:"dryrun"`
	Lines     int    `name:"lines" short:"n" usage:"head count"`
	Loop      bool   `name:"loop" usage:"loop playlist"`
	Reload    bool   `name:"reload" usage:"reload index"`
	Normalize bool   `name:"normalize" usage:"reload normalized index"`
	Shuffle   bool   `name:"shuffle" default:"true" usage:"shuffle music"`
	MusicRoot string `name:"music_root" short:"r" usage:"required, root directory of music files"`
	Debug     bool   `name:"debug" usage:"enable debug logs"`
	Quiet     bool   `name:"quiet" short:"q" usage:"quiet logs"`
	PlayOnly  bool   `name:"play" short:"s" usage:"read music file names from stdin instead of query, music_root is not required, options other than mpv, loop and dry are ignored"`
	Window    bool   `name:"window" short:"w" usage:"pretend GUI application"`
	//
	// external commands
	//
	MetafindCmd string `name:"metafind" default:"mf" usage:"metafind command, recommended: v0.6.1"`
	MpvCmd      string `name:"mpv" default:"mpv" usage:"mpv command, recommended: v0.40.0"`
	GrepCmd     string `name:"grep" default:"grep" usage:"grep command"`
	JqCmd       string `name:"jq" default:"jq" usage:"jq command, recommended: 1.8.1"`
	FfprobeCmd  string `name:"ffprobe" default:"ffprobe" usage:"ffprobe command, recommended: 7.1.1"`
	//
	// dependent variables
	//
	CacheDir string   `name:"-"`
	ListArgs []string `name:"-"`
	GrepArgs []string `name:"-"`
	//
	// etc
	//
	Writer io.Writer `name:"-" json:"-"`
}

func (c *Config) IndexFile() string {
	return filepath.Join(c.CacheDir, "index.json")
}

func (c *Config) NormalizedIndexFile() string {
	return filepath.Join(c.CacheDir, "normalized.index.json")
}
