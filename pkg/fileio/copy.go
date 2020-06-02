package fileio

import (
	"fmt"

	sys "golang.org/x/sys/unix"
)

func Copy(src, dest string, bufSize int) {
	inputFd := Open("src", src, sys.O_RDONLY, 0)
	defer Close(inputFd, "src", src)

	flagsCreateOrOverwrite := sys.O_CREAT | sys.O_WRONLY | sys.O_TRUNC

	outputFd := Open("dest", dest, flagsCreateOrOverwrite, PermsAllReadWrite)
	defer Close(outputFd, "dest", dest)

	buf := make([]byte, bufSize)

	for {
		numRead := ReadBuf(inputFd, buf, "src", src)
		if numRead == 0 {
			break
		}

		numWrite := WriteBuf(outputFd, buf[0:numRead], "dest", dest)
		if numWrite != numRead {
			panic(fmt.Sprintf("Partial write occurred dest file '%s' numRead=%d numWrite=%d", dest, numRead, numWrite))
		}
	}
}
