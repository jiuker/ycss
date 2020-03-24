package log

import (
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

func Log(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("%s:%d\r\n", file, line)
		if viper.GetBool("debug") {
			fmt.Println(v...)
		}
	}
}
