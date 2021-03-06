package bigbro

import (
	"fmt"
	"io"
	"sync"
)

// CSVFormatter is a formatter that formats in comma separated format.
type CSVFormatter struct {
	w  io.Writer
	mu sync.Mutex
}

// Format in comma separated file.
func (l CSVFormatter) Format(e Event) string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%d,%d,%s\n", e.Time.String(), e.Actor, e.Method, e.Target, e.Name, e.ID, e.Location, e.X, e.Y, e.ScreenWidth, e.ScreenHeight, e.Comment)
}

// Write the csv line to file.
func (l CSVFormatter) Write(e Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.w.Write([]byte(l.Format(e)))
	return err
}
