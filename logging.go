package bigbro

import (
	"io"
	"os"
	"time"
)

// Event is something which has happened on a web page.
type Event struct {
	// The element which has been triggered.
	Target string `json:"target",csv:"target"`
	// The name attribute of the element.
	Name string `json:"name",csv:"name"`
	// The id attribute of the element.
	ID string `json:"id",csv:"id"`
	// The method which triggered the event.
	Method string `json:"method",csv:"method"`
	// The web page location on the server.
	Location string `json:"location",csv:"location"`
	// Any additional information that can be useful.
	Comment string `json:"comment",csv:"comment"`
	// X position of the event.
	X int `json:"x",csv:"x"`
	// Y position of the event.
	Y int `json:"y",csv:"y"`
	// Width of the actors screen.
	ScreenWidth int `json:"screenWidth",csv:"screenWidth"`
	// Height of the actors screen.
	ScreenHeight int `json:"screenHeight",csv:"screenHeight"`
	// The time the Event happened.
	Time time.Time `json:"time",csv:"time"`
	// The actor that caused the Event.
	Actor string `json:"actor",csv:"actor"`
}

// Logger is the way in which logs are written to a file.
type Logger struct {
	f         io.Writer
	formatter Formatter
}

// NewLogger creates a new logger.
func NewLogger(name string, formatter Formatter) (Logger, error) {
	lf, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return Logger{}, err
	}
	lf.Truncate(0)

	return Logger{
		f:         lf,
		formatter: formatter,
	}, nil
}

// Log writes an event to the log file using the specified formatter.
func (l Logger) Log(e Event) error {
	line := l.formatter.Format(e)
	_, err := l.f.Write([]byte(line))
	return err
}
