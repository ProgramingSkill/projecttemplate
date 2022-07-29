package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	loadcfg()
	setDebugLevel(Config.Log.Level)
	functionInit()
	go signalHandle()
	startServer()
}
