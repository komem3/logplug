package logplug_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/komem3/logplug"
)

func TestLevelHook(t *testing.T) {
	for _, tt := range []struct {
		name   string
		option logplug.LevelConfig
		want   string
	}{
		{
			name:   "[DBG]ignore",
			option: logplug.LevelConfig{Levels: []logplug.Level{"DBG", "INFO"}, Min: "INFO"},
			want:   "",
		},
		{
			name:   "[DBG] remove blank",
			option: logplug.LevelConfig{Levels: []logplug.Level{"DBG", "INFO"}},
			want:   `{"level":"DBG","message":"remove blank"}`,
		},
		{
			name:   "info level apply",
			option: logplug.LevelConfig{Levels: []logplug.Level{"DBG", "INFO"}, Default: "INFO"},
			want:   `{"level":"INFO","message":"info level apply"}`,
		},
		{
			name:   "[WARN]change level field",
			option: logplug.LevelConfig{Levels: []logplug.Level{"DBG", "WARN"}, Field: "severity"},
			want:   `{"message":"change level field","severity":"WARN"}`,
		},
		{
			name:   "[DBG]log alias",
			option: logplug.LevelConfig{Alias: logplug.LevelAlias{"DBG": "DEBUG"}},
			want:   `{"level":"DEBUG","message":"log alias"}`,
		},
		{
			name:   "[INFO]log no alias",
			option: logplug.LevelConfig{Alias: logplug.LevelAlias{"DBG": "DEBUG"}},
			want:   `{"level":"INFO","message":"log no alias"}`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			log.New(logplug.NewJSONPlug(&buf, logplug.Hooks(
				logplug.LevelHook(tt.option),
			)), "", 0).Print(tt.name)

			if strings.TrimRight(buf.String(), "\n") != tt.want {
				t.Errorf("mismatch output\ngot:  %swant: %s", buf.String(), tt.want)
			}
		})
	}
}
