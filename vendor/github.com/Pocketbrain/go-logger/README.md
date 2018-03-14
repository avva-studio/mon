# go-logger

Logger package for Pocket Media Go Applications.

## About go-logger
The default logger for pm's go packages. This package is a small wrapper around
logrus (https://github.com/sirupsen/logrus) with the addition of setting up some
default options when running locally (debug purposes) or for other environments.

Also a hook is added for retrieving the filename, function and the line of where
the log entry call has been made.

## Example

### Usage
The logger can be used directly from the package like the following:
```go
import "github.com/Pocketbrain/go-logger"

func main() {
	plog.Debug("debug message")
	plog.Info("info message")
	plog.Warn("warn message")
	plog.Error("error message")
}
```
The default formatter is JSON and The default level is Warn.

There are several options available for the logger usage directly:
```go
import "github.com/Pocketbrain/go-logger"

func main() {
	plog.SetJSONFormatter() // To set the log formatter to JSON
	plog.SetTextFormatter() // To set the log formatter to Text

	plog.SetLevel(DebugLevel) // To set the log level to debug
	plog.SetLevel(InfoLevel) // To set the log level to info
	plog.SetLevel(WarnLevel) // To set the log level to warn
	plog.SetLevel(ErrorLevel) // To set the log level to error

	plog.SetDebugOptions() // To enable the debug options for the logger
}
```

### Instantiate
When using the package like following:

```go
var l plog.Logger

func init() {
	l = plog.New()

	// Instantiate with options
	l = plog.New(
		plog.JSONFormat(), // set the formatter to JSON
		plog.LevelInfo(), // set the level to level INFO
		plog.OutputStdOut(), // set the output to stdout
	)

	// Instantiate with debug options
	l = plog.New(DebugOptions())
}

func main() {
	l.Debug("debug message")
	l.Info("info message")
	l.Warning("warning message")
	l.Error("error message")

	l.with(logrus.Fields{
		"extra": "extra fields msg"
	}).Info("info with extra fields msg")

	l.SetLevel(DebugLevel) // Set the instantiated logger Level to debug
	l.SetLevel(InfoLevel) // Set the instantiated logger Level to info
	l.SetLevel(WarnLevel) // Set the instantiated logger Level to warn
	l.SetLevel(ErrorLevel) // Set the instantiated logger Level to error

	l.SetJSONFormatter() // To set the log formatter to JSON
	l.SetTextFormatter() // To set the log formatter to Text

	l.SetDebugOptions() // Enable the debug option for the instantiated logger
}
```

The following information should be retrieved:

```go
DEBU[2017-04-25T19:02:52+02:00] debug       file=README.md function=main line=22
INFO[2017-04-25T19:02:52+02:00] info        file=README.md function=main line=23
WARN[2017-04-25T19:02:52+02:00] warning     file=README.md function=main line=24
ERRO[2017-04-25T19:02:52+02:00] error       file=README.md function=main line=25
```

## Debug Mode
The debug mode allows for easy debugging. It can be set via flags. Logger
provides function StartDebugMode that accepts \*pflag.FlagSet as a parameter
and returns an error.

Flags:
```
debug - boolean, required
formatter - string, optional, allowed values: text|json, default: json
stdoutput - string, optional, allowed values: err|out, default: err
```

If debug flag is not set, or invalid value for formatter or stdoutput is passed
StartDebugMode returns an error.
When debug flag is set to true then logger will use level Debug and formatter
and stdoutput flag values and will ignore any changes of these options
further in the code other than a new call to StartDebugMode.

### Debug Mode Usage Example

Create flags:
```
var cmdRoot = &cobra.Command{
	Use:   "logger-example",
	Short: "short description",
	Long: "long description",
	SilenceUsage: true,
}

// Execute command
func Execute() {
	// if we want to implement debug flags in various commands we can add them as
	// persistent flags to the root command, otherwise we add them as normal
	// flags to a specific command
	cmdRoot.PersistentFlags().Bool("debug", false, "Defines if plog is used in debug mode")
	cmdRoot.PersistentFlags().String("formatter", "json", "Defines formatter if plog is used in debug mode")
	cmdRoot.PersistentFlags().String("stdoutput", "err", "Defines stdoutput if plog is used in debug mode")
	if err := cmdRoot.Execute(); err != nil {
		plog.Fatal(err)
		os.Exit(-1)
	}
}
```

Pass flags to the logger:
```
var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "short description",
	Long: "long description",
	Run: func(cmd *cobra.Command, args []string) {
		// pass cobra flags to the logger in order to be able to use logger in a debug mode
		err := plog.StartDebugMode(cmd.Flags())
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
```

Now we will be able to run the command using logger debug mode flags:
```
go run main.go serve --debug
go run main.go serve --debug --formatter=text
go run main.go serve --debug --stdoutput=out
```
