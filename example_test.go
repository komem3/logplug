package logplug_test

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/komem3/logplug"
)

func Example_globalLogger() {
	log.SetOutput(logplug.NewJSONPlug(os.Stderr,
		logplug.LogFlag(log.LstdFlags),
		logplug.Hooks(
			logplug.LevelHook(logplug.LevelConfig{
				Levels: []logplug.Level{"DBG", "INFO", "WARN", "ERR"},
				Min:    "INFO",
			}),
		)))

	log.Printf("[INFO]output test")
}

func Example_localLogger() {
	l := log.New(logplug.NewJSONPlug(os.Stderr,
		logplug.Hooks(
			logplug.LevelHook(logplug.LevelConfig{
				Levels: []logplug.Level{"DBUG", "INFO", "WARNING", "ERROR"},
				Alias: logplug.LevelAlias{
					"DBG":  "DBUG",
					"INFO": "INFO",
					"WARN": "WARNING",
					"ERR":  "ERROR",
				},
				Field: "severity",
			}),
		),
		logplug.LogFlag(log.Ldate|log.Lmicroseconds|log.Llongfile),
	), "[trace:1000][span:20]", log.Ldate|log.Lmicroseconds|log.Llongfile)

	l.Print("[ERR] output test")
}

type csvEncoder struct {
	writer *csv.Writer
}

func (c *csvEncoder) Encode(p *logplug.Plug, m *logplug.MessageElement) error {
	return c.writer.Write([]string{m.GetString("level"), m.GetString(p.MessageField())})
}

func ExampleEncoder() {
	csvwriter := csv.NewWriter(os.Stdout)
	defer csvwriter.Flush()

	err := csvwriter.Write([]string{"level", "message"})
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(logplug.NewPlug(&csvEncoder{writer: csvwriter}, logplug.Hooks(
		logplug.LevelHook(logplug.LevelConfig{
			Levels: []string{"DBG", "INFO", "ERR"},
			Field:  "level",
		}),
	)), "", 0)

	l.Print("[INFO]csv output")
	// output:
	// level,message
	// INFO,csv output
}

func timeLayoutHook(layout string) logplug.Hook {
	return func(enc logplug.Encoder) logplug.Encoder {
		return logplug.EncoderFunc(func(p *logplug.Plug, m *logplug.MessageElement) error {
			t := m.GetTime(p.TimestampField())
			if t.IsZero() {
				return enc.Encode(p, m)
			}

			m.Set(p.TimestampField(), t.Format(layout))
			return enc.Encode(p, m)
		})
	}
}

func ExampleHook() {
	const layout = "2006/01/02"
	l := log.New(logplug.NewJSONPlug(os.Stdout,
		logplug.LogFlag(log.Ldate),
		logplug.Hooks(timeLayoutHook(layout)),
	), "", log.Ldate)

	l.Print("time")
	// out: {"message":"time","timestamp":"2006/01/02"}
}
