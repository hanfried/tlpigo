package main

import (
	"github.com/alecthomas/kong"

	"github.com/hanfried/tlpigo/fileio"
)

func main() {
	var cli struct {
		BufSize int    `default:"1024"`
		Src     string `arg:""`
		Dst     string `arg:""`
	}

	kong.Parse(&cli)

	fileio.Copy(cli.Src, cli.Dst, cli.BufSize)
}
