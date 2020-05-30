package errors

import (
	"C"
	"os"
)

const EXIT_FAILURE = 1
const EXIT_WITH_DEFERS = true
const EXIT_WITHOUT_DEFERS = false

func Terminate(callDefers bool) {
	// Dump core if env variable GOTRACEBACK=crash is set
	if callDefers {
		// panic will call all the defers and also print a stack trace
		// not exactly what C exit(EXIT_FAILURE) does, but at least similar
		panic(EXIT_FAILURE)
	}
	// that's what C's _exit() function does
	os.Exit(EXIT_FAILURE)

	// XXX: Note, C's exit() is the recommended way to stop C programs
	// with flushing output buffers and calling atexit functions before exiting
	// where panic is the way for Go programs in exit case
	// While the pretty much same named function with an underscore in front _exit()
	// would just hard terminate the program and return the exit code without any graceful stop
}
