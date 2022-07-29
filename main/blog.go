package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var debuglevel = 4

const (
	//FatalLevel fatalerror
	FatalLevel = 0
	//ErrorLevel error happen
	ErrorLevel = 1
	//WarnLevel just warn something wrong
	WarnLevel = 2
	//InfoLevel info what happen
	InfoLevel = 3
	//DebugLevel for debug
	DebugLevel = 4
)

var (
	defaultLogLevel = DebugLevel
	errLevels       = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
)

var blog = log.New(os.Stdout, "", 0)
var logFlag = "projecttemplate"

func output(level string, logFlag string, format string, v ...interface{}) {
	pc, file, line, _ := runtime.Caller(2)
	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	f := runtime.FuncForPC(pc)
	fn := f.Name()

	for i := len(fn) - 1; i > 0; i-- {
		if fn[i] == '.' {
			fn = fn[i+1:]
			break
		}
	}

	if format == "" {
		blog.Printf("%v|%v|%v|%v|%v()|%v|%v", now(), level, logFlag, short, fn, line, fmt.Sprintln(v...))
	} else {
		blog.Printf("%v|%v|%v|%v|%v()|%v|%v", now(), level, logFlag, short, fn, line, fmt.Sprintf(format, v...))
	}

}
func now() string {
	t := time.Now()
	return fmt.Sprintf("%v-%02d-%v %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

//Debug for debug lowest level
func Debug(v ...interface{}) {
	if debuglevel >= DebugLevel {
		output("DEBUG", logFlag, "", v...)
	}
}

//Debugln same as Debug
func Debugln(v ...interface{}) {
	if debuglevel >= DebugLevel {
		output("DEBUG", logFlag, "", v...)
	}
}

//Debugf same as Debug
func Debugf(format string, v ...interface{}) {
	if debuglevel >= DebugLevel {
		output("DEBUG", logFlag, format, v...)
	}
}

//Info info level
func Info(v ...interface{}) {
	if debuglevel >= InfoLevel {
		output("INFO", logFlag, "", v...)
	}
}

//Infof info level
func Infof(format string, v ...interface{}) {
	if debuglevel >= InfoLevel {
		output("INFO", logFlag, format, v...)
	}
}

//Warning warning level
func Warning(v ...interface{}) {
	if debuglevel >= WarnLevel {
		output("WARN", logFlag, "", v...)
	}
}

//Warningln warning level
func Warningln(v ...interface{}) {
	if debuglevel >= WarnLevel {
		output("WARN", logFlag, "", v...)
	}
}

//Warningf warning level
func Warningf(format string, v ...interface{}) {
	if debuglevel >= WarnLevel {
		output("WARN", logFlag, format, v...)
	}
}

//Error error level
func Error(v ...interface{}) {
	if debuglevel >= ErrorLevel {
		output("ERROR", logFlag, "", v...)
	}
}

//Errorln error level
func Errorln(v ...interface{}) {
	if debuglevel >= ErrorLevel {
		output("ERROR", logFlag, "", v...)
	}
}

//Errorf error level
func Errorf(format string, v ...interface{}) {
	if debuglevel >= ErrorLevel {
		output("ERROR", logFlag, format, v...)
	}
}

//Fatal fatal error happen
func Fatal(v ...interface{}) {
	if debuglevel >= FatalLevel {
		output("FATAL", logFlag, "", v...)
	}
}

//Fatalf same as fatal
func Fatalf(format string, v ...interface{}) {
	if debuglevel >= FatalLevel {
		output("FATAL", logFlag, format, v...)
	}
}

//Logger init log func
func Logger(filename string, level string) *os.File {
	var f *os.File
	if filename != "" {
		var err error
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		blog = log.New(f, "", 0)
	}

	for k, v := range errLevels {
		if v == level {
			debuglevel = k
			Info("debuglevel:", v)
		}
	}

	return f
}

func setDebugLevel(level string) {
	for k, v := range errLevels {
		if v == level {
			debuglevel = k
			Fatal("debuglevel:", v)
		}
	}
}
