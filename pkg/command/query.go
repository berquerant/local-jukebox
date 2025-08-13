package command

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"

	"github.com/berquerant/local-jukebox/pkg/iox"
)

func NewQuery(
	metafind, grep string,
	query io.Reader,
	normalizedIndex string,
	listArgs, grepArgs []string,
	headCount int,
	shuffle bool,
	dest string,
) *Query {
	return &Query{
		metafind:        metafind,
		grep:            grep,
		query:           query,
		normalizedIndex: normalizedIndex,
		listArgs:        listArgs,
		grepArgs:        grepArgs,
		headCount:       headCount,
		shuffle:         shuffle,
		dest:            dest,
	}
}

// Query queries music from the index, grep, shuffle and head.
type Query struct {
	metafind        string
	grep            string
	query           io.Reader
	normalizedIndex string
	listArgs        []string
	grepArgs        []string
	headCount       int
	shuffle         bool
	dest            string
}

var _ Cmd = &Query{}

func (Query) Name() string { return "query" }

func (c *Query) Run(ctx context.Context) error {
	query, err := c.readQuery()
	if err != nil {
		return fmt.Errorf("%w: failed read query", err)
	}
	defer query.Close()

	dest, err := os.Create(c.dest)
	if err != nil {
		return fmt.Errorf("%w: failed to create dest", err)
	}
	defer dest.Close()

	listCmd := c.listCommand(ctx, query.Name())
	listCmd.Env = os.Environ()
	listStdout, err := listCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("%w: failed to create pipe for listing", err)
	}
	listCmd.Stderr = os.Stderr
	grepCmd := c.grepCommand(ctx)
	grepCmd.Env = os.Environ()
	grepCmd.Stdin = listStdout
	grepStdout, err := grepCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("%w: failed to create pipe for grep", err)
	}
	grepCmd.Stderr = os.Stderr

	logCmd(listCmd)
	if err := listCmd.Start(); err != nil {
		return fmt.Errorf("%w: failed to start listing", err)
	}
	logCmd(grepCmd)
	if err := grepCmd.Start(); err != nil {
		return fmt.Errorf("%w: failed to start grep", err)
	}

	var exitErr error
	resultLines, err := iox.ReadAllLines(grepStdout)
	if err != nil {
		exitErr = errors.Join(exitErr, fmt.Errorf("%w: failed to scan result", err))
	}
	if err := listCmd.Wait(); err != nil {
		exitErr = errors.Join(exitErr, fmt.Errorf("%w: failed to wait listing", err))
	}
	if err := grepCmd.Wait(); err != nil {
		if !isExitWith(err, 1) { // 1 means that no lines were selected
			exitErr = errors.Join(exitErr, fmt.Errorf("%w: failed to wait grep", err))
		}
	}
	if exitErr != nil {
		return exitErr
	}

	slog.Info("query", slog.Int("filtered", len(resultLines)))
	resultLines = c.head(resultLines)
	fmt.Fprintln(dest, strings.Join(resultLines, "\n"))
	slog.Info("query", slog.Int("passed", len(resultLines)))
	return nil
}

func (c *Query) listCommand(ctx context.Context, queryFile string) *exec.Cmd {
	arg := []string{
		"-i", c.normalizedIndex,
		"-e", "@" + queryFile,
	}
	arg = append(arg, c.listArgs...)
	return exec.CommandContext(ctx, c.metafind, arg...)
}

func (c *Query) grepCommand(ctx context.Context) *exec.Cmd {
	arg := c.grepArgs
	if len(arg) == 0 {
		arg = []string{".*"} // pass through
	}
	return exec.CommandContext(ctx, c.grep, arg...)
}

func (c *Query) readQuery() (*iox.TmpFile, error) {
	var buf []string
	scanner := bufio.NewScanner(c.query)
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// transform expressions of expr lang
	//
	// name matches "xxx"
	// name matches "yyy" and attr == "val"
	// ->
	// (name matches "xxx") or (name matches "yyy" and attr == "val")
	r := bytes.NewBufferString(strings.Join(buf, " or "))
	return iox.NewTmpFile(r, "query.txt")
}

func (c *Query) head(list []string) []string {
	if c.shuffle {
		rand.Shuffle(len(list), func(i, j int) {
			list[i], list[j] = list[j], list[i]
		})
	}
	if c.headCount < 1 || len(list) < c.headCount {
		return list
	}
	return list[:c.headCount]
}
