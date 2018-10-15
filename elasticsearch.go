package bigbro

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"time"
)

// LogstashEvent is an event that is recognised by Elasticsearch.
type LogstashEvent struct {
	Message   string    `json:"message"`
	Version   string    `json:"@version"`
	Timestamp time.Time `json:"@timestamp"`
	Type      string    `json:"type"`
	Host      string    `json:"host"`
	Event     Event     `json:"event"`
}

// ElasticsearchFormatter is a log formatter that can output to Elasticsearch.
type ElasticsearchFormatter struct {
	index   string
	version string
	client  *elastic.Client
}

func (f ElasticsearchFormatter) Format(e Event) string {
	b, _ := json.Marshal(f.transformEvent(e))
	return string(b)
}

func (f ElasticsearchFormatter) Write(e Event) error {
	ctx := context.Background()
	_, err := f.client.Index().
		Index(f.index).
		Type("event").
		BodyJson(f.transformEvent(e)).
		Do(ctx)
	return err
}

// transformEvent returns an Elasticsearch compatible version of an Event.
func (f ElasticsearchFormatter) transformEvent(e Event) LogstashEvent {
	return LogstashEvent{
		Message:   e.Name,
		Version:   f.version,
		Timestamp: e.Time,
		Type:      e.Method,
		Host:      e.Location,
		Event:     e,
	}
}

// NewElasticsearchFormatter creates a new formatter for Elasticsearch.
func NewElasticsearchFormatter(index, version, url string) (ElasticsearchFormatter, error) {
	c, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return ElasticsearchFormatter{}, err
	}
	return ElasticsearchFormatter{
		index:   index,
		version: version,
		client:  c,
	}, nil
}
