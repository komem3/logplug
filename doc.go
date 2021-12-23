/*
Package logplug implements a plug of standard logging package.
This package enables json output, level logging and so on to standard logger.

Plug implements io.Writer.
Therefore, Plug can be specified as the output destination of log.
Plug converts the output of log to the following format.
	log.SetFlags(log.Ldate | log.Lshortfile)
	log.Print("log message")
	// normal output: "[label:test] 2006-01-02 doc.go:1:log message""
	// convert:
	map[string]interface{}{
		"message": "logMessage",
		"timestamp": time.Date(2006,1,2,0,0,0,0,0,time.Local),
		"location": "doc.go:1",
	}

Pass the converted map to the encoder as an argument.
Therefore, you can change the output format of log by the encoder of Plug.
	log.SetOutput(log.NewJSONPlug(os.Stderr))
	log.Print("output")
	// output: {"message":"output"}

Use prefix if you want to use the log to include custom fields.
Custom fields set the values in [] as key and value.
	log.Printf("[key:value] custom")
	// convert:
	map[string]interface{}{
		"key": "value",
		"message": "custom",
	}

You can use SetPrefix to set fields that are common to all logs.
	log.SetPrefix("[key]:value")
	log.Printf("custom")
	// convert:
	map[string]interface{}{
		"key": "value",
		"message": "custom",
	}

It supports hooks before encoding.
Hooks allow level logging, field changes and so on.
	logplug.NewJSONPlug(os.Stderr, logplug.Hooks(
		logplug.LevelHook(logplug.LevelConfig{
			Levels: []logplug.Level{"DBG", "INFO", "WARN", "ERR"},
			Min:    "INFO",
		}),
	))
*/
package logplug
