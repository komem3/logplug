package benchmarks_test

import (
	"log"
	"os"
	"testing"

	"github.com/komem3/logplug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkLogOutputWithoutArgs(b *testing.B) {
	tmp, err := os.CreateTemp("", "logfile*.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	l := log.New(tmp, "prefix", log.LstdFlags)

	b.ResetTimer()
	b.Run("standard logger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Print("Output")
		}
	})

	l = log.New(logplug.NewJSONPlug(l.Writer(), logplug.LogFlag(l.Flags())), l.Prefix(), l.Flags())
	b.ResetTimer()
	b.Run("standard logger with plug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Print("Output")
		}
	})

	zlog := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		tmp,
		zap.InfoLevel,
	))
	if err != nil {
		b.Fatal(err)
	}
	b.Run("zap logger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			zlog.Info("Output")
		}
	})
}
