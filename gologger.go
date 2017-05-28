package gologger

import (
	"fmt"
	"time"
	"os"
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

func (l *GoLogger)log(a ...interface{}) {
	args := fmt.Sprint(a)
	msg := time.Now().Format(time.UnixDate) + " :: " + args[1: len(args) - 1]
	fmt.Println(msg)

	// TODO: Add support for different log levels
	if l.LogLevel == 1 {
		go func(message string) {
			f, err:= os.OpenFile(l.LogPath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Error: could not write to octopus log file:", l.LogPath)
				return
			}
			defer f.Close()

			message += "\n"
			_, err = f.WriteString(message)
			if err != nil {
				fmt.Println("Error: failed to write message to log file:", l.LogPath)
			}
		}(msg)
	}
}
