package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jiuker/ycss/handle"

	"github.com/jiuker/ycss/cfg"
)

func main() {
	// parse base config file path
	var baseConfigPath string
	flag.StringVar(&baseConfigPath, "base", "./res/config", "set base config file path")
	cfg.SetBasePath(baseConfigPath)
	handle.StartHandle()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGSTOP)
	<-exit
}
