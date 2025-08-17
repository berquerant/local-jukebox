# local-jukebox

```
❯ jukebox --help
jukebox - play music files by querying your local library

# Usage

  jukebox [flags] [--] [grep args...] [-- mf -i additional args...]

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

External filter.

  jukebox -r /root/dir/of/music -x query.txt --dry -- -- -v | grep 'WORD' | jq -r .path | jukebox --play

# Prerequisites

- mf https://github.com/berquerant/metafind
- mpv https://github.com/mpv-player/mpv
- grep
- jq https://github.com/jqlang/jq
- ffprobe https://ffmpeg.org/ffprobe.html

# Flags
      --debug               enable debug logs
  -l, --dry                 dryrun
      --ffprobe string      ffprobe command, recommended: 7.1.1 (default "ffprobe")
      --grep string         grep command (default "grep")
      --jq string           jq command, recommended: 1.8.1 (default "jq")
  -n, --lines int           head count
      --loop                loop playlist
      --metafind string     metafind command, recommended: v0.6.1 (default "mf")
      --mpv string          mpv command, recommended: v0.40.0 (default "mpv")
  -r, --music_root string   required, root directory of music files
      --normalize           reload normalized index
  -s, --play                read music file names from stdin instead of query, music_root is not required, options other than mpv, loop and dry are ignored
  -x, --query string        music query (default "stdin")
  -q, --quiet               quiet logs
      --reload              reload index
      --shuffle             shuffle music (default true)
  -w, --window              pretend GUI application
```
