package main

import (
	"github.com/alecthomas/kong"

	"github.com/hanfried/tlpigo/fileio"
)

func main() {
	var cli struct {
		BufSize  int    `default:"1024"`
		HoleSize int64  `default:"0"`
		Src      string `arg:""`
		Dst      string `arg:""`
	}

	kong.Parse(&cli)

	if cli.HoleSize > 0 {
		fileio.SparseCopy(cli.Src, cli.Dst, cli.BufSize, cli.HoleSize)
	} else {
		fileio.Copy(cli.Src, cli.Dst, cli.BufSize)
	}
}
