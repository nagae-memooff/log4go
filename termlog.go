// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	rc       chan *LogRecord
	close_wg sync.WaitGroup
	closing  bool
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() *ConsoleLogWriter {
	records := &ConsoleLogWriter{
		rc: make(chan *LogRecord, LogBufferLength),
	}

	go records.run(stdout)
	return records
}

func (w *ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w.rc {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
		}
		fmt.Fprint(out, "[", timestr, "] [", levelStrings[rec.Level], "] ", rec.Message, "\n")
		w.close_wg.Done()
	}

}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	if w.closing {
		return
	}

	w.close_wg.Add(1)
	w.rc <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w *ConsoleLogWriter) Close() {
	w.closing = true
	w.close_wg.Wait()
	close(w.rc)
}
