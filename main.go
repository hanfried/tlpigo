package main

import (
	err "github.com/hanfried/tlpigo/lib/errors"
)

func main() {
	err.Terminate(err.EXIT_WITH_DEFERS)
}
