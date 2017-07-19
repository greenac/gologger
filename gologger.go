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

func (l *GoLogger)Log(a ...interface{}) {
	args := fmt.Sprint(a)
	msg := time.Now().Format(time.UnixDate) + " :: " + args[2: len(args) - 2]
	l.coloredOutput(outputNormal, msg)

	// TODO: Add support for different log levels
	if l.LogLevel == 1 {
		go func(message string) {
			f, err:= os.OpenFile(l.LogPath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
			if err != nil {
				l.coloredOutput(outputError, "Error: could not write to octopus log file:", l.LogPath)
				return
			}
			defer f.Close()

			message += "\n"
			_, err = f.WriteString(message)
			if err != nil {
				l.coloredOutput(outputError, "Error: failed to write message to log file:", l.LogPath)
			}
		}(msg)
	}
}
