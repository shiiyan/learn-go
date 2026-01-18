package main

import "log"

func myLog(format string, args ...any) {
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}
