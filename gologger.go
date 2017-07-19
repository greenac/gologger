package gologger

import (
	"time"
	"os"
	"github.com/fatih/color"
	"fmt"
)

type outputType int
const (
	outputError outputType = 1000
	outputNormal outputType = 1001
)

type GoLogger struct {
	LogLevel int
	LogPath string
	timeFormat string
	isSetup bool
}

// TODO: Add support for different log levels
func (l *GoLogger)Setup() {
	if !l.isSetup {
		if l.timeFormat == "" {
			l.timeFormat = time.UnixDate
		}

		l.isSetup = true
	}
}

func (l *GoLogger)coloredOutput(ot outputType, a... interface{}) {
	var c *color.Color
	if ot == outputError {
		c = color.New(color.FgRed).Add(color.Bold)
	} else {
		c = color.New(color.FgCyan)
	}

	c.Println(a...)
}

func (l *GoLogger)writeToFile(message string) {
	if l.LogLevel == 1 {
		go func(msg string) {
			f, err:= os.OpenFile(l.LogPath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
			if err != nil {
				l.coloredOutput(outputError, "Error: could not write to octopus log file:", l.LogPath)
				return
			}
			defer f.Close()

			msg += "\n"
			_, err = f.WriteString(msg)
			if err != nil {
				l.coloredOutput(outputError, "Error: failed to write message to log file:", l.LogPath)
			}
		}(message)
	}
}

func (l *GoLogger)log(ot outputType, a ...interface{}) {
	args := fmt.Sprint(a)
	msg := time.Now().Format(time.UnixDate) + " :: " + args[3: len(args) - 3]
	l.coloredOutput(ot, msg)
	l.writeToFile(msg)
}

func (l *GoLogger)Log(a ...interface{}) {
	l.log(outputNormal, a)
}

func (l *GoLogger)Error(a ...interface{}) {
	l.Log(outputError, a)
}
