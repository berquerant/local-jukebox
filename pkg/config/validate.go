package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (c *Config) Validate() error {
	if err := c.validateExecutable(); err != nil {
		return err
	}
	if err := c.validateCacheDir(); err != nil {
		return err
	}
	if err := c.validateMusicRoot(); err != nil {
		return err
	}
	return nil
}

func (c *Config) validateMusicRoot() error {
	if x, err := os.Stat(c.MusicRoot); err == nil {
		if x.IsDir() {
			return nil
		}
	}
	return errors.New("music_root is not a directory")
}

func (c *Config) validateCacheDir() error {
	x, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("%w: user cache dir is missing", err)
	}
	c.CacheDir = filepath.Join(x, "local-jukebox")
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil {
		return fmt.Errorf("%w: cannot create user cache dir %s", err, c.CacheDir)
	}
	return nil
}

func (c *Config) validateExecutable() error {
	list := []string{
		c.MetafindCmd,
		c.MpvCmd,
		c.GrepCmd,
		c.JqCmd,
		c.FfprobeCmd,
	}
	for _, x := range list {
		if err := isExecutable(x); err != nil {
			return err
		}
	}
	return nil
}

func isExecutable(v string) error {
	if _, err := exec.LookPath(v); err != nil {
		return fmt.Errorf("%w: %s is not executable", err, v)
	}
	return nil
}
