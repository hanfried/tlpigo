package fileio

import (
	"fmt"

	sys "golang.org/x/sys/unix"
)

func Tee(dest string, append bool, bufSize int) {
	flags := map[bool]int{
		false: sys.O_CREAT | sys.O_WRONLY | sys.O_TRUNC,
		true:  sys.O_CREAT | sys.O_WRONLY | sys.O_APPEND,
	}[append]

	outputFd := Open("dest", dest, flags, PermsAllReadWrite)
	defer Close(outputFd, "dest", dest)

	buf := make([]byte, bufSize)

	for {
		numRead := ReadBuf(sys.Stdin, buf, "input", "STDIN")
		if numRead == 0 {
			break
		}

		numWriteDest := WriteBuf(outputFd, buf[0:numRead], "dest", dest)
		if numWriteDest != numRead {
			panic(fmt.Sprintf("Partial write occurred dest file '%s' numRead=%d numWrite=%d", dest, numRead, numWriteDest))
		}

		numWriteStdout := WriteBuf(sys.Stdout, buf[0:numRead], "output", "STDOUT")
		if numWriteStdout != numRead {
			panic(fmt.Sprintf("Partial write occurred when writing to STDOUT numRead=%d numWrite=%d", numRead, numWriteStdout))
		}
	}
}
