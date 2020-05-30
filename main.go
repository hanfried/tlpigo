package main

import (
	err "github.com/hanfried/tlpigo/errors"
)

func main() {
	err.Terminate(err.EXIT_WITH_DEFERS)
}
