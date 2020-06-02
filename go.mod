module github.com/hanfried/tlpigo

go 1.14

replace (
	github.com/hanfried/tlpigo/errors => ./internal/errors
	github.com/hanfried/tlpigo/fileio => ./pkg/fileio
)

require (
	github.com/alecthomas/kong v0.2.9
	github.com/hanfried/tlpigo/errors v0.0.0-00010101000000-000000000000
	github.com/hanfried/tlpigo/fileio v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20200523222454-059865788121
)
