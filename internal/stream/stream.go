// Package stream provides primitives for handling the reading of a JSON formatted file
// which can be a big file in size.
// It is not possible to think to load into memory serveral Gb of data, expecially when the
// memory is limited
// due to not so much time, the implementation has been roughly inspired by
// http://www.inanzzz.com/index.php/post/5ehg/reading-and-decoding-a-large-json-file-as-in-streaming-fashion-with-golang
package stream

import (
	"encoding/json"
	"fmt"
	"microservice-exercise/internal/data_model"
	"os"
)

// StreamService is an interface for loading PortItem as a JSON stream file
type StreamService interface {
	Load(path string)
	Watch() <-chan data_model.PortItem
}

// Stream represents a JSON stream, from which we read a file which can be a very big file.
// The stream is related to the data_model we have to handle, which main entity is the PortItem
type Stream struct {
	stream chan data_model.PortItem
}

// NewStream return new JSON stream in order to read the PortItem(s)
func NewStream() *Stream {
	return &Stream{
		stream: make(chan data_model.PortItem),
	}
}

// Load reads JSON file in stream for decoding the PortItem(s)
func (s Stream) Load(path string) {
	file, err := os.Open(path)
	if err != nil {
		s.stream <- data_model.PortItem{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	t, err := decoder.Token()
	if err != nil {
		s.stream <- data_model.PortItem{Error: fmt.Errorf("JSON token: %w", err)}
		return
	}
	if t != json.Delim('{') {
		s.stream <- data_model.PortItem{Error: fmt.Errorf("expected {, got %v", t)}
		return
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			s.stream <- data_model.PortItem{Error: fmt.Errorf("JSON token: %w", err)}
			return
		}
		ID := t.(string)

		var port data_model.Port
		if err := decoder.Decode(&port); err != nil {
			s.stream <- data_model.PortItem{Error: fmt.Errorf("decode: %w", err)}
			return
		}

		s.stream <- data_model.PortItem{ID: ID, Port: port}
	}
}

// Watch watches for stream entries
func (s Stream) Watch() <-chan data_model.PortItem {
	return s.stream
}
