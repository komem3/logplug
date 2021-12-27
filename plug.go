package logplug

import (
	"log"
	"strings"
	"sync"
	"time"
)

// Encoder defines a encode for each field in the log.
type Encoder interface {
	Encode(p *Plug, m *MessageElement) error
}

// EncoderFunc is adapter of Encoder.
type EncoderFunc func(p *Plug, m *MessageElement) error

// Encode implements Encoder.
func (f EncoderFunc) Encode(p *Plug, m *MessageElement) error {
	return f(p, m)
}

// Hook is a hook of encoder.
type Hook func(enc Encoder) Encoder

// MessageElement store elements of message.
type MessageElement struct {
	elements map[string]interface{}
}

var messageElementPool = sync.Pool{
	New: func() interface{} {
		return &MessageElement{
			elements: make(map[string]interface{}, 3),
		}
	},
}

func (m *MessageElement) GetString(key string) string {
	v, _ := m.elements[key].(string)
	return v
}

// Set set v to key of elements.
// This method override exist value.
// If you want to set string value, recommend to use AddString.
func (m *MessageElement) Set(key string, v interface{}) {
	m.elements[key] = v
}

// AddString add v to key of elements.
func (m *MessageElement) AddString(key string, v string) {
	if str, ok := m.elements[key].(string); ok {
		m.elements[key] = str + v
	} else {
		m.elements[key] = v
	}
}

func (m *MessageElement) GetBool(key string) bool {
	v, _ := m.elements[key].(bool)
	return v
}

func (m *MessageElement) GetTime(key string) time.Time {
	v, _ := m.elements[key].(time.Time)
	return v
}

func (m *MessageElement) Elements() map[string]interface{} {
	return m.elements
}

// Plug is standard log plug.
type Plug struct {
	encoder Encoder
	hooks   []Hook

	messageField   string
	timeStampField string
	locationField  string
	flag           int
}

// NewPlug create new log plug.
func NewPlug(encoder Encoder, opts ...Option) *Plug {
	p := &Plug{
		messageField:   "message",
		timeStampField: "timestamp",
		locationField:  "location",
	}

	for _, opt := range opts {
		opt(p)
	}

	for i := len(p.hooks) - 1; i >= 0; i-- {
		encoder = p.hooks[i](encoder)
	}
	p.encoder = encoder

	return p
}

// Write implements io.Writer.
func (p *Plug) Write(msgb []byte) (n int, err error) {
	mel := messageElementPool.Get().(*MessageElement)
	msg := string(msgb)

	// log flag process
	if p.flag&log.Lmsgprefix != 0 {
		t, index := p.extractTimestamp(msg)
		if !t.IsZero() {
			msg = msg[index+1:]
			mel.Set(p.timeStampField, t)
		}

		f := p.extractFile(msg)
		if f != "" {
			msg = msg[len(f)+2:]
			mel.Set(p.locationField, f)
		}
	}

	// prefix process
	for len(msg) > 0 && msg[0] == '[' {
		end := strings.IndexRune(msg, ':')
		if end == -1 {
			break
		}

		prefixEnd := strings.Index(msg[end:], "]")
		if prefixEnd == -1 {
			break
		}

		value := msg[end+1 : prefixEnd+end]
		if b, ok := p.parseBool(value); ok {
			mel.Set(msg[1:end], b)
		} else {
			mel.AddString(msg[1:end], value)
		}

		msg = msg[prefixEnd+end+1:]
	}

	// log flag process
	if mel.GetTime(p.timeStampField).IsZero() && mel.GetString(p.locationField) == "" {
		t, index := p.extractTimestamp(msg)
		if !t.IsZero() {
			msg = msg[index+1:]
			mel.Set(p.timeStampField, t)
		}
		f := p.extractFile(msg)
		if f != "" {
			msg = msg[len(f)+2:]
			mel.Set(p.locationField, f)
		}
	}

	// prefix process
	// pattern: "[prefix]timestamp[prefix]message"
	for len(msg) > 0 && msg[0] == '[' {
		end := strings.IndexRune(msg, ':')
		if end == -1 {
			break
		}

		prefixEnd := strings.Index(msg[end:], "]")
		if prefixEnd == -1 {
			break
		}

		value := msg[end+1 : prefixEnd+end]
		if b, ok := p.parseBool(value); ok {
			mel.Set(msg[1:end], b)
		} else {
			mel.AddString(msg[1:end], value)
		}
		msg = msg[prefixEnd+end+1:]
	}

	msg = strings.TrimLeft(strings.TrimRight(msg, "\n"), " ")
	mel.AddString(p.messageField, msg)

	if err := p.encoder.Encode(p, mel); err != nil {
		return 0, err
	}
	for key := range mel.elements {
		delete(mel.elements, key)
	}
	messageElementPool.Put(mel)
	return len(msgb), nil
}

func (p *Plug) MessageField() string {
	return p.messageField
}

func (p *Plug) TimestampField() string {
	return p.timeStampField
}

func (p *Plug) LocationField() string {
	return p.locationField
}

func (p *Plug) LogFlag() int {
	return p.flag
}

func (*Plug) secondIndex(s string, r rune) int {
	var hit bool
	for i, sr := range s {
		match := sr == r
		if hit && match {
			return i
		}
		hit = hit || match
	}
	return -1
}

func (p *Plug) extractTimestamp(msg string) (time.Time, int) {
	if p.flag&log.Ldate != 0 {
		if p.flag&log.Lmicroseconds != 0 {
			date := msg[:p.secondIndex(msg, ' ')]
			timestamp, err := time.Parse("2006/01/02 15:04:05.999999", date)
			if err != nil {
				return time.Time{}, 0
			}
			return timestamp, len(date)
		} else if p.flag&log.Ltime != 0 {
			date := msg[:p.secondIndex(msg, ' ')]
			timestamp, err := time.Parse("2006/01/02 15:04:05", date)
			if err != nil {
				return time.Time{}, 0
			}
			return timestamp, len(date)
		} else {
			date := msg[:strings.IndexRune(msg, ' ')]
			timestamp, err := time.Parse("2006/01/02", date)
			if err != nil {
				return time.Time{}, 0
			}
			return timestamp, len(date)
		}
	}
	return time.Time{}, 0
}

func (p *Plug) extractFile(msg string) string {
	if p.flag&(log.Lshortfile|log.Llongfile) != 0 {
		index := p.secondIndex(msg, ':')
		if index <= 0 {
			return ""
		}
		return msg[:index]
	}
	return ""
}

func (p *Plug) parseBool(msg string) (b bool, ok bool) {
	switch msg {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}
