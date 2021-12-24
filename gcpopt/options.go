package gcpopt

import (
	"log"
	"strings"

	"github.com/komem3/logplug"
)

// LogFlags is expected log flag.
// should set this value to log.
//	log.SetFlags(gcpopt.LogFlags)
var LogFlags = log.Llongfile

const levelField = "severity"

// DefaultLevelConfig is a config according to LogSeverity of gcp.
var DefaultLevelConfig = logplug.LevelConfig{
	Levels:  []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"},
	Default: "INFO",
	Alias: map[string]string{
		"DBG":  "DEBUG",
		"WARN": "WARNING",
		"ERR":  "ERROR",
	},
	Field: levelField,
}

// ErrorReportHook converts the specified levels to the format of an error report.
//	https://cloud.google.com/error-reporting/docs/formatting-error-messages#@type
func ErrorReportHook(levels ...logplug.Level) logplug.Hook {
	alartLevel := make(map[logplug.Level]bool, len(levels))
	for _, level := range levels {
		alartLevel[level] = true
	}
	return func(enc logplug.Encoder) logplug.Encoder {
		return logplug.EncoderFunc(func(p *logplug.Plug, m *logplug.MessageElement) error {
			if _, ok := alartLevel[m.GetString(levelField)]; ok {
				m.AddString("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent")
			}
			return enc.Encode(p, m)
		})
	}
}

// LocationModifyHook convert the value of location in log to sourceLocation.
func LocationModifyHook() logplug.Hook {
	return func(enc logplug.Encoder) logplug.Encoder {
		return logplug.EncoderFunc(func(p *logplug.Plug, m *logplug.MessageElement) error {
			location := m.GetString(p.LocationField())

			index := strings.LastIndexByte(location, ':')
			if index == -1 && index-1 == len(location) {
				return enc.Encode(p, m)
			}
			sourceLocation := struct {
				File string `json:"file"`
				Line string `json:"line"`
			}{
				File: location[:index],
				Line: location[index+1:],
			}

			delete(m.Elements(), p.LocationField())

			m.Set("logging.googleapis.com/sourceLocation", sourceLocation)
			return enc.Encode(p, m)
		})
	}
}

// NewGCPOptions provides options for gcp logging.
// conf will modify gcpopt.DefaultLevelConfig.
//
// This assumes that gcpopt.logFlags has been set for log.
//	log.SetFlags(gcpopt.LogFlags)
//
func NewGCPOptions(minLevel logplug.Level) []logplug.Option {
	conf := DefaultLevelConfig
	conf.Min = minLevel
	return []logplug.Option{
		logplug.LogFlag(LogFlags),
		logplug.Hooks(
			logplug.LevelHook(conf),
			ErrorReportHook("ALERT", "EMERGENCY"),
			LocationModifyHook(),
		),
	}
}
