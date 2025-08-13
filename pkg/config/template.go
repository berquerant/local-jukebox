package config

import (
	"bytes"
	"text/template"
)

const metafindConfigTemplate = `
root:
  - "$MUSIC_ROOT"
probe:
  - |
    {{ .FfprobeCmd }} -v error -hide_banner -show_entries format -of json=c=1 @ARG | {{ .JqCmd }} .format.tags -c
pname:
  - ffp
`

func (c *Config) MetafindConfig() (string, error) {
	t := template.Must(template.New("metafind").Parse(metafindConfigTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, c); err != nil {
		return "", err
	}
	return buf.String(), nil
}
