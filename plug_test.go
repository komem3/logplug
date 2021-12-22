package logplug_test

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/komem3/logplug"
)

type plugTestCase struct {
	name   string
	prefix string
	flag   int
	want   string
}

func TestJSONPlug(t *testing.T) {
	for _, tt := range []plugTestCase{
		{
			name: "text only",
			want: `{"message":"text only"}`,
		},
		{
			name: "with prefix", prefix: "[trace:1000]",
			want: `{"message":"with prefix","trace":"1000"}`,
		},
		{
			name: "prefix is boolean", prefix: "[trace:true]",
			want: `{"message":"prefix is boolean","trace":true}`,
		},
		{
			name: "with multi prefix", prefix: "[trace1:1000][trace2:2000]",
			want: `{"message":"with multi prefix","trace1":"1000","trace2":"2000"}`,
		},
		{
			name: "append field value", prefix: "[trace:one][trace:two]",
			want: `{"message":"append field value","trace":"onetwo"}`,
		},
		{
			name: "[prefix:before]message prefix", prefix: "[trace:one]",
			want: `{"message":"message prefix","prefix":"before","trace":"one"}`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			log.New(logplug.NewJSONPlug(&buf), tt.prefix, tt.flag).
				Print(tt.name)

			if strings.TrimRight(buf.String(), "\n") != tt.want {
				t.Errorf("mismatch output\ngot:  %swant: %s", buf.String(), tt.want)
			}
		})
	}
}

func TestJSONPlug_WithLogFlag(t *testing.T) {
	for _, tt := range []plugTestCase{
		{
			name: "with date", flag: log.Ldate,
			want: `{"message":"with date","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T00:00:00Z"}`,
		},
		{
			name: "with date+time", flag: log.Ldate | log.Ltime,
			want: `{"message":"with date\+time","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z"}`,
		},
		{
			name: "with date+milisec", flag: log.Ldate | log.Lmicroseconds,
			want: `{"message":"with date\+milisec","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{0,6}Z"}`,
		},
		{
			name: "with date+time+milisec", flag: log.Ldate | log.Ltime | log.Lmicroseconds,
			want: `{"message":"with date\+time\+milisec","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{0,6}Z"}`,
		},
		{
			name: "with short file", flag: log.Lshortfile,
			want: `{"location":"plug_test\.go:[0-9]+","message":"with short file"}`,
		},
		{
			name: "with long file", flag: log.Llongfile,
			want: `{"location":"[[:graph:]]+plug_test\.go:[0-9]+","message":"with long file"}`,
		},
		{
			name: "with prefix", flag: log.Ldate, prefix: "[label:test]",
			want: `{"label":"test","message":"with prefix","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T00:00:00Z"}`,
		},
		{
			name: " with msg: prefix", flag: log.Ldate | log.Lmsgprefix | log.Lshortfile, prefix: "[label:test]",
			want: `{"label":"test","location":"plug_test\.go:[0-9]+","message":"with msg: prefix","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T00:00:00Z"}`,
		},
		{
			name: "[with:prefix]msg prefix", flag: log.Ldate, prefix: "[label:test]",
			want: `{"label":"test","message":"msg prefix","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T00:00:00Z","with":"prefix"}`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			log.New(logplug.NewJSONPlug(&buf, logplug.LogFlag(tt.flag)), tt.prefix, tt.flag).
				Print(tt.name)

			match, err := regexp.Match(tt.want, buf.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Errorf("mismatch output\ngot:  %swant: %s", buf.String(), tt.want)
			}
		})
	}
}
