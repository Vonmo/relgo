/*
 * Created: 2017.11
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

var Level = 1
var StackLvl = 1
var ShowFilePos = true

func Init(level string, filepos bool, dest string) {
	log.SetFlags(log.LUTC | log.Lmicroseconds)
	if dest != "stdout" {
		f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println("rewrite logs to file")
	}
	ShowFilePos = filepos
	switch level {
	case "debug":
		Level = 1
	case "info":
		Level = 2
	case "error":
		Level = 3
	default:
		Level = 1
	}
}

func Debug(v ...interface{}) {
	if Level > 1 {
		return
	}
	format := "[debug] %s"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	if Level > 1 {
		return
	}
	format = "[debug] " + format
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, v...)
}

func Info(v ...interface{}) {
	if Level > 2 {
		return
	}
	format := "[info] %s"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	if Level > 2 {
		return
	}
	format = "[info] " + format
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, v...)
}

func Error(v ...interface{}) {
	if Level > 3 {
		return
	}
	format := "[error] %s"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	if Level > 2 {
		return
	}
	format = "[error] " + format
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, v...)
}

func Panic(v ...interface{}) {
	format := "[panic] %s"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Panicf(format, fmt.Sprint(v...))
}

func Panicf(format string, v ...interface{}) {
	format = "[panic] " + format
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Panicf(format, v...)
}

func Fatal(v ...interface{}) {
	format := "[fatal] %s"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Fatalf(format, fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) {
	format = "[fatal] " + format
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Fatalf(format, v...)
}

func Print(v ...interface{}) {
	if Level > 2 {
		return
	}
	format := "[info] %s"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	if Level > 2 {
		return
	}
	format = "[info] " + format
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, fmt.Sprint(v...))
}

func Println(v ...interface{}) {
	if Level > 2 {
		return
	}
	format := "[info] %s\n"
	if ShowFilePos {
		_, file, line, ok := runtime.Caller(StackLvl)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		format = file + ":" + strconv.Itoa(line) + ":" + format
	}

	log.Printf(format, fmt.Sprint(v...))
}
