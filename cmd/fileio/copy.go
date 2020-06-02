package main

import (
	"fmt"

	"github.com/alecthomas/kong"

	"github.com/hanfried/tlpigo/fileio"
	sys "golang.org/x/sys/unix"
)

func main() {
	var cli struct {
		BufSize int    `default:"1024"`
		Src     string `arg:""`
		Dst     string `arg:""`
	}

	kong.Parse(&cli)

	inputFd := fileio.Open("src", cli.Src, sys.O_RDONLY, 0)
	defer fileio.Close(inputFd, "src", cli.Src)

	flagsCreateOrOverwrite := sys.O_CREAT | sys.O_WRONLY | sys.O_TRUNC

	outputFd := fileio.Open("dest", cli.Dst, flagsCreateOrOverwrite, fileio.PermsAllReadWrite)
	defer fileio.Close(outputFd, "dest", cli.Dst)

	buf := make([]byte, cli.BufSize)

	for {
		numRead := fileio.ReadBuf(inputFd, buf, "src", cli.Src)
		if numRead == 0 {
			break
		}

		numWrite := fileio.WriteBuf(outputFd, buf[0:numRead], "dest", cli.Dst)
		if numWrite != numRead {
			panic(fmt.Sprintf("Partial write occurred dest file '%s' numRead=%d numWrite=%d", cli.Dst, numRead, numWrite))
		}
	}
}
