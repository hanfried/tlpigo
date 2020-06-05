package main

import (
	"github.com/alecthomas/kong"

	"github.com/hanfried/tlpigo/fileio"
)

func main() {
	var cli struct {
		Append  bool   `default:"false"`
		BufSize int    `default:"1024"`
		Dst     string `arg:""`
	}

	kong.Parse(&cli)

	fileio.Tee(cli.Dst, cli.Append, cli.BufSize)
}
