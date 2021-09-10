package main

import (
	"github.com/widaT/logbus"
)

func main() {
	logbus.SetLogLevel("default", logbus.TraceLevel)
	logbus.Infof("%s,%s", "aaaaaaaaaa", "9999")
	logbus.Debugf("%s", "77777777777777")
	logbus.Tracef("%s", "77777777777777")
	log := logbus.NewLogger(logbus.DebugLevel, "main")
	log.Debugf("aaaaaaaaaaa")
	logs := logbus.GetLoggers()
	for _, l := range logs {
		l.Debugf("aaaaaaaaa")
	}
}
