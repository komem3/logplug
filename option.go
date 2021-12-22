package logplug

// Option is option of Plug.
type Option func(p *Plug)

// LogFlag set flag of log.
// This value should match the value of flag of log.
//
// see available flag:
// 	go doc log.Ldate
func LogFlag(flag int) Option {
	return func(p *Plug) {
		p.flag = flag
	}
}

// Hooks set hooks of Plug.
func Hooks(hooks ...Hook) Option {
	return func(p *Plug) {
		p.hooks = append(p.hooks, hooks...)
	}
}

// MessageFiled set field name of message.
func MessageFiled(filed string) Option {
	return func(p *Plug) {
		p.messageField = filed
	}
}

// TimestampField set field name of timestamp.
func TimestampField(field string) Option {
	return func(p *Plug) {
		p.timeStampField = field
	}
}

// LocationField set field name of location.
func LocationField(field string) Option {
	return func(p *Plug) {
		p.locationField = field
	}
}
