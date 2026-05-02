package logger

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

type PrefixWriter struct {
	Prefix      string
	Dest        io.Writer
	isNewLine   bool
	mu          sync.Mutex
	lastWrite   time.Time
}

func NewPrefixWriter(prefix string, dest io.Writer) *PrefixWriter {
	return &PrefixWriter{
		Prefix:    prefix,
		Dest:      dest,
		isNewLine: true,
		lastWrite: time.Now(),
	}
}

func (w *PrefixWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastWrite = time.Now()

	var buf bytes.Buffer
	for _, b := range p {
		if w.isNewLine {
			buf.WriteString(w.Prefix)
			w.isNewLine = false
		}
		buf.WriteByte(b)
		if b == '\n' {
			w.isNewLine = true
		}
	}

	// Write to console
	fmt.Print(buf.String())

	// Forward to destination
	return w.Dest.Write(p)
}

func (w *PrefixWriter) GetLastWriteTime() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastWrite
}
