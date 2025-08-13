package iox

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

func WriteFile(r io.Reader, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

type TmpFile struct {
	path string
}

func NewTmpFile(r io.Reader, path string) (*TmpFile, error) {
	p, err := os.MkdirTemp("", "local-jukebox")
	if err != nil {
		return nil, err
	}
	t := &TmpFile{
		path: filepath.Join(p, path),
	}
	if err := WriteFile(r, t.path); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TmpFile) Name() string            { return t.path }
func (t *TmpFile) Open() (*os.File, error) { return os.Open(t.path) }
func (f *TmpFile) Close() error            { return os.RemoveAll(filepath.Dir(f.path)) }

func ReadAllLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func ExistFile(path string) bool {
	if x, err := os.Stat(path); err == nil {
		return !x.IsDir()
	}
	return false
}
