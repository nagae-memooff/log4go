// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	rc   chan *LogRecord
	quit chan bool
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() ConsoleLogWriter {
	records := ConsoleLogWriter{
		rc:   make(chan *LogRecord, LogBufferLength),
		quit: make(chan bool),
	}

	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w.rc {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
		}
		fmt.Fprint(out, "[", timestr, "] [", levelStrings[rec.Level], "] ", rec.Message, "\n")
	}

	w.quit <- true
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w.rc <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w.rc)
	<-w.quit
}
