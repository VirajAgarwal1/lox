package error_pak

import (
	"fmt"
	"runtime"
)

var runtimme_lineNum int
var runtime_file string
var runtime_pc uintptr

func RetErr(msg string, stacked_err error) error {
	runtime_pc, runtime_file, runtimme_lineNum, _ = runtime.Caller(1)
	if stacked_err != nil {
		if len(msg) != 0 {
			return fmt.Errorf("%d| %v/%v  ->  %v\n%w", runtimme_lineNum-2, runtime_file, runtime.FuncForPC(runtime_pc).Name(), msg, stacked_err)
		}
		return fmt.Errorf("%d| %v/%v\n%w", runtimme_lineNum-2, runtime_file, runtime.FuncForPC(runtime_pc).Name(), stacked_err)
	}
	return fmt.Errorf("%d| %v/%v)\n%v", runtimme_lineNum-2, runtime_file, runtime.FuncForPC(runtime_pc).Name(), msg)
}
