package events

import (
	"bytes"
	"context"
	"io"

	"github.com/golang/protobuf/jsonpb"
)

// Marshaler provides a default JSON marshaler.
var Marshaler = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	Indent:       "\t",
}

// Marshal marshals the event to JSON.
func Marshal(ev *Event) (string, error) {
	return Marshaler.MarshalToString(ev)
}

// Unmarshal unmarshals the event from JSON.
func Unmarshal(b []byte, ev *Event) error {
	r := bytes.NewReader(b)
	return jsonpb.Unmarshal(r, ev)
}

var CompactMarshaler = &jsonpb.Marshaler{}

type JSONWriter struct {
	Output io.Writer
}

func (j *JSONWriter) WriteEvent(ctx context.Context, ev *Event) error {
	err := CompactMarshaler.Marshal(j.Output, ev)
	if err != nil {
		return err
	}
	_, err = io.WriteString(j.Output, "\n")
	return err
}
