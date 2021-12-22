# logplug

This package enables json output or level logging to standard logger.

[![Go Reference](https://pkg.go.dev/badge/github.com/komem3/logplug.svg)](https://pkg.go.dev/github.com/komem3/logplug)

## Usage

```go
log.SetOutput(logplug.NewJSONPlug(os.Stderr,
	logplug.LogFlag(log.LstdFlags),
	logplug.Hooks(
		logplug.LevelHook(logplug.LevelConfig{
			Levels: []logplug.Level{"DBG", "INFO", "WARN", "ERR"},
			Min:    "INFO",
		}),
	)))

log.Printf("[INFO]output test")
// output: {"level":"INFO","message":"output test"}
```

## Options

- [Options for GCP](./gcpopt)

## Examples

- [GCP Logger](./gcpopt/example)

## License

MIT

