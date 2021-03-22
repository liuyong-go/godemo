package ygo

import (
	"fmt"
	"runtime"
)

func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return err
			}
		}
		return nil
	}
}
func try(fn func() error, cleaner func()) error {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(2)
			//初始化过日志就打印日志，未初始化就输出
			fmt.Println(err, file, line)
		}
	}()
	return fn()
}
