package command

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func NewNormalize(dest, src string) *Normalize {
	return &Normalize{
		src:  src,
		dest: dest,
	}
}

// Normalize normalizes the metafind index as unicode strings.
type Normalize struct {
	src  string
	dest string
}

var _ Cmd = &Normalize{}

func (Normalize) Name() string { return "normalize" }

func (c *Normalize) Run(ctx context.Context) error {
	src, err := os.Open(c.src)
	if err != nil {
		return fmt.Errorf("%w: failed to open src", err)
	}
	defer src.Close()

	dest, err := os.Create(c.dest)
	if err != nil {
		return fmt.Errorf("%w: failed to create dest", err)
	}
	defer dest.Close()

	type P struct {
		Path string `json:"path"`
	}

	slog.Info("normalize", slog.String("src", c.src), slog.String("dest", c.dest))
	var (
		count  int
		failed int
	)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		count++
		x := scanner.Text()
		var p P
		if err := json.Unmarshal([]byte(x), &p); err != nil {
			slog.Warn("normalize", slog.Any("error", err))
			failed++
			continue
		}
		n := normalizeString(x)
		s := fmt.Sprintf(`{"path":"%s","n":%s}`, p.Path, n)
		fmt.Fprintln(dest, s)
	}
	slog.Info("normalize", slog.Int("count", count), slog.Int("failed", failed))
	return scanner.Err()
}
