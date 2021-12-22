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

// DefaultLevelConfig is a config according to LogSeverity of gcp.
var DefaultLevelConfig = logplug.LevelConfig{
	Levels:  []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"},
	Default: "INFO",
	Min:     "DEBUG",
	Alias: map[string]string{
		"DBG":  "DEBUG",
		"WARN": "WARNING",
		"ERR":  "ERROR",
	},
	Field: "severity",
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
//
// This assumes that gcpopt.logFlags has been set for log.
//	log.SetFlags(gcpopt.LogFlags)
//
// conf will modify gcpopt.DefaultLevelConfig.
//	gcpopt.NewGCPOptions(func(conf *logplug.LevelConfig) {
//		conf.Min = "INFO"
//	})
// If conf is nil, it will use gcpopt.DefaultLevelConfig.
func NewGCPOptions(modifyConfig func(conf *logplug.LevelConfig)) []logplug.Option {
	conf := DefaultLevelConfig
	if modifyConfig != nil {
		modifyConfig(&conf)
	}
	return []logplug.Option{
		logplug.LogFlag(LogFlags),
		logplug.Hooks(
			logplug.LevelHook(conf),
			LocationModifyHook(),
		),
	}
}
