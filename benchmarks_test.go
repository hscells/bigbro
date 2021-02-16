package bigbro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkCSVLogging(b *testing.B) {
	logger, err := NewCSVLogger(fmt.Sprintf("/tmp/bb_%s.log", time.Now().Format(time.RFC3339)))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We create a random event.
		event := Event{
			Target:       "window",
			Name:         "",
			ID:           "",
			Method:       "click",
			Location:     "https://example.com",
			Comment:      "",
			X:            rand.Int(),
			Y:            rand.Int(),
			ScreenWidth:  1000,
			ScreenHeight: 1000,
			Time:         time.Now(),
			Actor:        "A1",
		}

		// Encode and decode the event to simulate it coming over the wire.
		var j bytes.Buffer
		err := json.NewEncoder(&j).Encode(event)
		if err != nil {
			b.Fatal(err)
		}
		var sentEvent Event
		err = json.NewDecoder(&j).Decode(&sentEvent)
		if err != nil {
			b.Fatal(err)
		}

		b.SetBytes(int64(len(logger.formatter.Format(sentEvent))))

		// Finally, log the event.
		err = logger.Log(sentEvent)
		if err != nil {
			b.Fatal(err)
		}
	}
}
