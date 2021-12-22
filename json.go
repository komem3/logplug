package logplug

import (
	"encoding/json"
	"io"
)

// jsonEncoder wrap Encoder for json encoder.
type jsonEncoder struct {
	encoder *json.Encoder
}

// Encode implements Encoder.
func (i *jsonEncoder) Encode(_ *Plug, m *MessageElement) error {
	return i.encoder.Encode(m.Elements())
}

// NewJSONPlug create a plug that converts log to json.
func NewJSONPlug(w io.Writer, opts ...Option) *Plug {
	return NewPlug(&jsonEncoder{json.NewEncoder(w)}, opts...)
}
