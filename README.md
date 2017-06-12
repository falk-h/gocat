# gocat

A blazing fast [lolcat](https://github.com/busyloop/lolcat) replacement.

[![demo](https://asciinema.org/a/4bfe41lnu96yx1rvgxrop6g5v.png)](https://asciinema.org/a/4bfe41lnu96yx1rvgxrop6g5v?autoplay=1)

## Building

Building requires a functioning [go installation](https://golang.org/doc/install#install).

- Clone the repository
- `go build gocat.go`

## Usage

`gocat [OPTION]... [FILE]...`

- `-a, --animate` — Animate the output
- `-d, --duration=<d>` — Animation duration (default: 12)
    - Animation duration is specified in number of frames per line
- `-f, --force` — Force color output
    - By default, gocat does not rainbowify messages if its output is redirected
    - Useful for writing a colourful message to `/etc/motd`, etc
- `-F, --freq=<f>` — Rainbow frequency (default: 2)
    - Defines how much the color changes for each character
- `-i, --invert` — Invert the output
- `-n, --number` — Number all output lines
- `-O, --offset=<o>` — Vertical offset (default: 2)
    - Same as `-F`, but for lines.
- `-S, --seed=<s>` — RNG seed, 0 means random (default: 0)
- `-s, --speed=<s>` — Animation speed (default: 20)
    - Animation speed is specified in frames per second
- `-h, --help` — Display the help text
- `--version` — Display version information
