package fileio

import (
	"fmt"

	sys "golang.org/x/sys/unix"
)

const (
	PermsUserReadWrite  = uint32(sys.S_IRUSR | sys.S_IWUSR)
	PermsGroupReadWrite = uint32(sys.S_IRGRP | sys.S_IWGRP)
	PermsOtherReadWrite = uint32(sys.S_IROTH | sys.S_IWOTH)
	PermsAllReadWrite   = PermsUserReadWrite | PermsGroupReadWrite | PermsOtherReadWrite
)

func Open(desc, fname string, mode int, perms uint32) int {
	fd, err := sys.Open(fname, mode, perms)
	if err != nil {
		panic(fmt.Sprintf("Can't open %s file '%s': %s", desc, fname, err.Error()))
	}

	return fd
}

func Close(fd int, desc, fname string) {
	err := sys.Close(fd)
	if err != nil {
		panic(fmt.Sprintf("Can't close %s file '%s' (fd %d): %s", desc, fname, fd, err.Error()))
	}
}

func ReadBuf(fd int, buf []byte, desc, fname string) int {
	numRead, err := sys.Read(fd, buf)
	if err != nil {
		panic(fmt.Sprintf("Error reading %s file '%s': %s", desc, fname, err.Error()))
	} else if numRead < 0 {
		panic(fmt.Sprintf("Should not happen: numRead=%d < 0 (%s file '%s')", numRead, desc, fname))
	}

	return numRead
}

func WriteBuf(fd int, buf []byte, desc, fname string) int {
	numWrite, err := sys.Write(fd, buf)
	if err != nil {
		panic(fmt.Sprintf("Error writing %s file '%s': %s", desc, fname, err.Error()))
	}

	return numWrite
}
