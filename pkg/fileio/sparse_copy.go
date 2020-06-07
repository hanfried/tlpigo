package fileio

import (
	"fmt"
	"io"

	sys "golang.org/x/sys/unix"
)

type HoleSearchResult struct {
	length               int64
	endPos, bufReadCount int
	endBuf               []byte
	eof                  bool
}

func findEndOfZeros(startPos, bufSize int, bufStart []byte, inputFd int, src string) (res HoleSearchResult) {
	res.length = 1
	res.endPos = startPos + 1
	res.endBuf = bufStart

	var readBuf []byte
outer:
	for {
		for res.endPos < len(res.endBuf) {
			if res.endBuf[res.endPos] != 0 {
				break outer
			} else {
				res.length++
				res.endPos++
			}
		}

		if readBuf == nil {
			readBuf = make([]byte, bufSize)
		}
		numRead := ReadBuf(inputFd, readBuf, "src", src)
		res.bufReadCount++
		res.endBuf = readBuf[0:numRead]
		res.endPos = 0

		if numRead == 0 {
			res.eof = true
			break
		}
	}

	return res
}

func seekHole(fd int, size int64) {
	_, err := sys.Seek(fd, size, io.SeekEnd)
	if err != nil {
		panic(fmt.Sprintf("lseek error '%s'", err.Error()))
	}
}

type SingletonEmptyBuffer struct {
	buf  []byte
	size int
}

func (singleton *SingletonEmptyBuffer) WriteNTimes(fd int, fname string, n int) {
	if singleton.buf == nil {
		singleton.buf = make([]byte, singleton.size)
	}

	for _i := 0; _i < n; _i++ {
		WriteBuf(fd, singleton.buf, "dest", fname)
	}
}

func SparseCopy(src, dest string, bufSize int, holeSize int64) {
	inputFd := Open("src", src, sys.O_RDONLY, 0)
	defer Close(inputFd, "src", src)

	flagsCreateOrOverwrite := sys.O_CREAT | sys.O_WRONLY | sys.O_TRUNC

	outputFd := Open("dest", dest, flagsCreateOrOverwrite, PermsAllReadWrite)
	defer Close(outputFd, "dest", dest)

	var buf []byte

	emptyBuf := SingletonEmptyBuffer{size: bufSize}
	readBuf := make([]byte, bufSize)
loop_through_file:
	for {
		numRead := ReadBuf(inputFd, readBuf, "src", src)
		if numRead == 0 { // end of file
			break loop_through_file
		}
		buf = readBuf[0:numRead]

		pos := 0
	loop_through_buffer:
		for pos < len(buf) {
			if buf[pos] != 0 {
				pos++
				continue loop_through_buffer
			}

			zeros := findEndOfZeros(pos, bufSize, buf, inputFd, src) // we are looking for ASCII code 0
			if !zeros.eof && zeros.length >= holeSize {
				// Seeking at end of file without writing anything behind does not change the file
				// so from this point of view, we can't regard NULs at end of file as a hole
				// so have to write the NULs hard without seeking
				// XXX: Discuss whether it might be an idea to write zeros.length - 1 NULs and then write 1 NUL
				WriteBuf(outputFd, buf[0:pos], "dest", dest)
				seekHole(outputFd, zeros.length)
				buf = zeros.endBuf[zeros.endPos:]
				pos = 0
			} else {
				if zeros.bufReadCount > 0 {
					WriteBuf(outputFd, buf, "dest", dest)
					emptyBuf.WriteNTimes(outputFd, dest, zeros.bufReadCount-1)
				}
				pos = zeros.endPos
				buf = zeros.endBuf
			}
		}
		WriteBuf(outputFd, buf, "dest", dest)
	}
}
