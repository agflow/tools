package agerr

import (
	"fmt"
	"log"
	"time"
)

// CallAndLog calls function which may return error and logs it.
// The intention of this function is to be used with `go` and `defer` clauses.
func CallAndLog(f func() error) {
	Log(f())
}

// Log logs error unless nil
func Log(err error) {
	if err != nil {
		nowStr := time.Now().Format("2006/01/02 15:04:05")
		msg := fmt.Sprintf("unhandled error %+v", err)
		log.Printf("%s%s %s %s%s", "\033[31m", nowStr, "[ERROR]", msg, "\033[0m")
	}
}

// Assert panics and breaks the program unless err is nil
func Assert(err error) {
	if err != nil {
		log.Panic(err)
	}
}
