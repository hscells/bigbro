package bigbro

import "fmt"

// Formatter is a way of specifying how to format an event for logging.
type Formatter interface {
	Format(e Event) string
}

// CSVFormatter is a formatter that formats in comma separated format.
type CSVFormatter struct{}

// Format in comma separated file.
func (l CSVFormatter) Format(e Event) string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%d,%d,%s\n", e.Time.String(), e.Actor.Identifier, e.Method, e.Target, e.Name, e.ID, e.Location, e.X, e.Y, e.ScreenWidth, e.ScreenHeight, e.Comment)
}
