package logplug

import (
	"strings"
)

// Level is logging level.
type Level = string

// LevelAlias is alias of level.
type LevelAlias map[string]Level

// LevelConfig is option of LevelEncoder.
type LevelConfig struct {
	Levels  []Level
	Default Level
	Min     Level
	Alias   LevelAlias
	Field   string
}

// LevelHook is hook of parse level and level filter.
// LevelHook use the string in the first [] of the Message field as the level.
//	log.Print("[DBG] message") // level=DBG
// If config.Default and config.Min are empty, set config.Levels[0].
//
// simple:
//
//	logplug.LevelHook(logplug.LevelConfig{Levels: []logplug.Level{"DBG", "WARN", "ERR"}})
//
// setting min level:
//
//	logplug.LevelHook(logplug.LevelConfig{Levels: []logplug.Level{"DBG","ERR"}, Min: "ERR"})
//
// level alias(DBG -> DEBUG):
//
//	logplug.LevelHook(logplug.LevelConfig{Alias: logplug.LevelAlias{"DBG": "DEBUG"}})
func LevelHook(config LevelConfig) Hook {
	badLevel := make(map[string]struct{}, len(config.Levels))

	if config.Default == "" && len(config.Levels) > 0 {
		config.Default = config.Levels[0]
	}
	if config.Min == "" && len(config.Levels) > 0 {
		config.Min = config.Levels[0]
	}
	if config.Field == "" {
		config.Field = "level"
	}

	for _, level := range config.Levels {
		if level == config.Min {
			break
		}
		badLevel[level] = struct{}{}
	}

	return func(enc Encoder) Encoder {
		return EncoderFunc(func(p *Plug, m *MessageElement) error {
			msg := m.GetString(p.MessageField())
			if len(msg) == 0 {
				return enc.Encode(p, m)
			}

			var (
				level Level
				end   = strings.IndexRune(msg, ']')
			)
			if end == -1 || end == len(msg)-1 {
				level = config.Default
			} else {
				level = msg[1:end]

				if msg[end+1] == ' ' {
					m.Set(p.MessageField(), strings.TrimLeft(msg[end+1:], " "))
				} else {
					m.Set(p.MessageField(), msg[end+1:])
				}
			}

			if config.Alias != nil {
				if l, ok := config.Alias[level]; ok {
					level = l
				}
			}

			if _, ok := badLevel[level]; ok {
				return nil
			}
			m.AddString(config.Field, level)
			return enc.Encode(p, m)
		})
	}
}
