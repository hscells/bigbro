package bigbro

// Formatter is a way of specifying how to format an event for logging.
type Formatter interface {
	Format(e Event) string
	Write(e Event) error
}