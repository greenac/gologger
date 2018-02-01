package gologger

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

type outputType int

const (
	OutputError   outputType = 1000
	OutputNormal  outputType = 1001
	OutputWarning outputType = 1002
)

type GoLogger struct {
	LogLevel   outputType
	LogPath    string
	timeFormat string
	isSetup    bool
}

// TODO: Add support for different log levels
func (l *GoLogger) Setup() {
	if !l.isSetup {
		if l.timeFormat == "" {
			l.timeFormat = time.UnixDate
		}

		l.isSetup = true
	}
}

func (l *GoLogger) coloredOutput(ot outputType, a ...interface{}) {
	var c *color.Color
	switch ot {
	case OutputError:
		c = color.New(color.FgRed).Add(color.Bold)
	case OutputWarning:
		c = color.New(color.FgYellow).Add(color.Bold)
	default:
		c = color.New(color.FgCyan)
	}

	c.Println(a...)
}

func (l *GoLogger) writeToFile(message string) {
	l.Setup()
	if l.LogPath == "" {
		return
	}

	if l.LogLevel == 1 {
		go func(msg string) {
			f, err := os.OpenFile(l.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				l.coloredOutput(OutputError, "Error: could not write to octopus log file:", l.LogPath)
				return
			}
			defer f.Close()

			msg += "\n"
			_, err = f.WriteString(msg)
			if err != nil {
				l.coloredOutput(OutputError, "Error: failed to write message to log file:", l.LogPath)
			}
		}(message)
	}
}

func (l *GoLogger) log(ot outputType, a ...interface{}) {
	l.Setup()
	args := fmt.Sprint(a)
	var pre string
	switch ot {
	case OutputError:
		pre = "ERROR: "
	case OutputWarning:
		pre = "WARNING: "
	case OutputNormal:
		pre = ""
	default:
		fmt.Println("Error: output type:", ot, "is unknown")
		pre = ""
	}

	start := -1
	end := -1
	for i := 0; i < len(args)/2; i++ {
		if start == -1 && args[i] != '[' {
			start = i
		}

		if end == -1 && args[len(args) - i - 1] != ']' {
			end = len(args) - i
		}

		if start != -1 && end != -1 {
			break
		}
	}

	if start == -1 {
		start = 0
	}

	if end == -1 {
		end = len(args)
	}

	msg := time.Now().Format(time.UnixDate) + " :: " + pre + args[start: end]
	l.coloredOutput(ot, msg)
	l.writeToFile(msg)
}

func (l *GoLogger) Log(a ...interface{}) {
	l.log(OutputNormal, a)
}

func (l *GoLogger) Error(a ...interface{}) {
	l.log(OutputError, a)
}

func (l *GoLogger) Warn(a ...interface{}) {
	l.log(OutputWarning, a)
}
