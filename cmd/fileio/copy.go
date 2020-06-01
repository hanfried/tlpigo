package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	sys "golang.org/x/sys/unix"
)

func openFile(desc, fname string, mode int, perms uint32) int {
	fd, err := sys.Open(fname, mode, perms)
	if err != nil {
		panic(fmt.Sprintf("Can't open %s file '%s': %s", desc, fname, err.Error()))
	}

	return fd
}

func closeFile(fd int, desc, fname string) {
	err := sys.Close(fd)
	if err != nil {
		panic(fmt.Sprintf("Can't close %s file '%s' (fd %d): %s", desc, fname, fd, err.Error()))
	}
}

func main() {
	var cli struct {
		BufSize int    `default:"1024"`
		Src     string `arg:""`
		Dst     string `arg:""`
	}

	kong.Parse(&cli)

	inputFd := openFile("src", cli.Src, sys.O_RDONLY, 0)
	defer closeFile(inputFd, "src", cli.Src)

	flagsCreateOrOverwrite := sys.O_CREAT | sys.O_WRONLY | sys.O_TRUNC
	permsUserReadWrite := sys.S_IRUSR | sys.S_IWUSR
	permsGroupReadWrite := sys.S_IRGRP | sys.S_IWGRP
	permsOtherReadWrite := sys.S_IROTH | sys.S_IWOTH
	permsAllReadWrite := uint32(permsUserReadWrite | permsGroupReadWrite | permsOtherReadWrite)

	outputFd := openFile("dest", cli.Dst, flagsCreateOrOverwrite, permsAllReadWrite)
	defer closeFile(outputFd, "dest", cli.Dst)

	buf := make([]byte, cli.BufSize)

	for {
		numRead, err := sys.Read(inputFd, buf)
		if err != nil {
			panic(fmt.Sprintf("Error reading src file '%s': %s", cli.Src, err.Error()))
		} else if numRead < 0 {
			panic(fmt.Sprintf("Should not happen: numRead=%d < 0", numRead))
		}

		if numRead == 0 {
			break
		}

		numWrite, err := sys.Write(outputFd, buf[0:numRead])
		if err != nil {
			panic(fmt.Sprintf("Error writing dest file '%s': %s", cli.Dst, err.Error()))
		}

		if numWrite != numRead {
			panic(fmt.Sprintf("Partial write occurred dest file '%s' numRead=%d numWrite=%d", cli.Dst, numRead, numWrite))
		}
	}
}
