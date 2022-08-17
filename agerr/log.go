package agerr

import "github.com/agflow/tools/log"

// CallAndLog calls function which may return error and logs it.
// The intention of this function is to be used with `go` and `defer` clauses.
func CallAndLog(f func() error) {
	Log(f())
}

// Log logs error unless nil
func Log(err error) {
	if err != nil {
		log.Errorf("unhandled error %+v", err)
	}
}

// Assert panics and breaks the program unless err is nil
func Assert(err error) {
	if err != nil {
		log.Panic(err)
	}
}
