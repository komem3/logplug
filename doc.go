/*
Package logplug implements a plug of standard logging package.
This package enables json output or level logging to standard logger.

Plug implements io.Writer.
Therefore, Plug can be specified as the output destination of log.
You can change the output format of log by the encoder of Plug.
	log.SetOutput(log.NewJSONPlug(os.Stderr))

It also supports hooks before encoding.
Hooks allow level logging, field changes and so on.
	logplug.NewJSONPlug(os.Stderr, logplug.Hooks(
		logplug.LevelHook(logplug.LevelConfig{
			Levels: []logplug.Level{"DBG", "INFO", "WARN", "ERR"},
			Min:    "INFO",
		}),
	))

*/
package logplug
