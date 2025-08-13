package iox_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/berquerant/local-jukebox/pkg/iox"
	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	const content = "CONTENT"
	path := filepath.Join(t.TempDir(), "testd", "test")
	assert.Nil(t, iox.WriteFile(bytes.NewBufferString(content), path))
	f, err := os.Open(path)
	if !assert.Nil(t, err) {
		return
	}
	defer f.Close()
	got, err := io.ReadAll(f)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, content, string(got))
}

func TestTmpFile(t *testing.T) {
	const (
		content = "CONTENT"
		path    = "testd/test"
	)
	tmpf, err := iox.NewTmpFile(bytes.NewBufferString(content), path)
	if !assert.Nil(t, err) {
		return
	}
	assert.True(t, strings.HasSuffix(tmpf.Name(), path))
	func() {
		f, err := tmpf.Open()
		if !assert.Nil(t, err) {
			return
		}
		defer f.Close()
		got, err := io.ReadAll(f)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, content, string(got))
	}()
	if !assert.Nil(t, tmpf.Close()) {
		return
	}
	_, err = os.Stat(tmpf.Name())
	assert.True(t, os.IsNotExist(err))
}

func TestReadAllLines(t *testing.T) {
	for _, tc := range []struct {
		title string
		input string
		want  []string
	}{
		{
			title: "3 lines",
			input: "l1\nl2\nl3\n",
			want:  []string{"l1", "l2", "l3"},
		},
		{
			title: "a line",
			input: "l1",
			want:  []string{"l1"},
		},
		{
			title: "empty",
			want:  nil,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := iox.ReadAllLines(bytes.NewBufferString(tc.input))
			if !assert.Nil(t, err) {
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestExistFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test")
	f, err := os.Create(path)
	if !assert.Nil(t, err) {
		return
	}
	f.Close()

	for _, tc := range []struct {
		title string
		path  string
		want  bool
	}{
		{
			title: "not exist",
			path:  path + "not_exist",
			want:  false,
		},
		{
			title: "exist but directory",
			path:  t.TempDir(),
			want:  false,
		},
		{
			title: "exist",
			path:  path,
			want:  true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, iox.ExistFile(tc.path))
		})
	}
}
