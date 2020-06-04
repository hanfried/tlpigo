package main

import (
	"fmt"
	"io"
	"strconv"

	sys "golang.org/x/sys/unix"

	"github.com/alecthomas/kong"

	"github.com/hanfried/tlpigo/fileio"
)

type ErrorParseSeekOperation struct {
	arg            string
	parsingProblem string
	fromErr        error
}

func (e ErrorParseSeekOperation) Error() string {
	desc := fmt.Sprintf("Could not parse arg '%s': %s", e.arg, e.parsingProblem)

	switch e.fromErr {
	case nil:
		return desc
	default:
		return fmt.Sprintf("%s: %s", desc, e.fromErr.Error())
	}
}

type seekOperation struct {
	cmd          byte
	len          int64  // for read commands 'r', 'R'
	offset       int64  // for seek command 's'
	writeContent string // for write command 'w'
}

func parseOperation(s string) (op seekOperation, err error) {
	op.cmd = s[0]
	switch op.cmd {
	case 'r', 'R':
		op.len, err = strconv.ParseInt(s[1:], 10, 64)
		if err != nil {
			err = ErrorParseSeekOperation{arg: s, parsingProblem: "length is not an int64", fromErr: err}
		}
	case 'w':
		op.writeContent = s[1:]
		if len(op.writeContent) == 0 {
			err = ErrorParseSeekOperation{arg: s, parsingProblem: "no write content"}
		}
	case 's':
		op.offset, err = strconv.ParseInt(s[1:], 10, 64)
		if err != nil {
			err = ErrorParseSeekOperation{arg: s, parsingProblem: "offset is not an int64:", fromErr: err}
		}
	default:
		err = ErrorParseSeekOperation{arg: s, parsingProblem: "unknown seek operation"}
	}

	return op, err
}

func main() {
	var cli struct {
		File       string   `arg:""`
		Operations []string `arg:"" help:"r<LEN>|R<LEN>|w<STRING>|s<OFFSET>"`
	}

	kong.Parse(&cli)

	flagsReadWriteOrCreate := sys.O_RDWR | sys.O_CREAT
	fd := fileio.Open("seek", cli.File, flagsReadWriteOrCreate, fileio.PermsAllReadWrite)
	ops := make([]seekOperation, len(cli.Operations))

	for i, opstr := range cli.Operations {
		op, err := parseOperation(opstr)
		if err != nil {
			fmt.Printf("seek: error: %s\n", err.Error())
			sys.Exit(1)
		}

		ops[i] = op
	}

	for i, op := range ops {
		switch op.cmd {
		case 'r', 'R': // display bytes at current offset as text (r) or as hex (R)
			buf := make([]byte, op.len)

			numRead := fileio.ReadBuf(fd, buf, "seek", cli.File)
			if numRead == 0 {
				fmt.Printf("end of seek file '%s' reached\n", cli.File)
			} else {
				fmt.Printf("%s: ", cli.Operations[i])
				charfmt := map[byte]string{
					'r': "%c",
					'R': "%02x ",
				}[op.cmd]
				for _, c := range buf[0:numRead] {
					if op.cmd == 'r' && !strconv.IsPrint(rune(c)) {
						c = '?'
					}
					fmt.Printf(charfmt, c)
				}
			}
		case 'w':
			numWritten := fileio.WriteBuf(fd, []byte(op.writeContent), "seek", cli.File)
			fmt.Printf("%c%s: write %d bytes\n", op.cmd, op.writeContent, numWritten)
		case 's':
			off, err := sys.Seek(fd, op.offset, io.SeekStart)
			if err != nil {
				panic(fmt.Sprintf("lseek error: '%s'", err.Error()))
			}

			fmt.Printf("%c%d: seek succeeded, set offset to %d\n", op.cmd, op.offset, off)
		}
	}
}
