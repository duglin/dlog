package dlog

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DLogger struct {
	sync.Mutex
	verbose    int
	indent     string
	log        *log.Logger
	autoindent bool // do we indent children log msgs with: |
	keywords   []string
}

var std = &DLogger{
	verbose:    1,
	log:        log.Default(),
	autoindent: false,
}

func (log *DLogger) String() string {
	return fmt.Sprintf("Log: verbose:%d indent:%q autoindent:%v keywords:%v",
		log.verbose, log.indent, log.autoindent, log.keywords)
}

func (log *DLogger) GetVerbose() int      { return log.verbose }
func (log *DLogger) SetVerbose(v int)     { log.verbose = v }
func (log *DLogger) SetAutoIndent(v bool) { log.autoindent = v }

// str is of the form xxx,xxx where xxx is either the version int
// for use with VPrint or a keyword for use with KPrint
func (log *DLogger) AddVerboseString(str string) {
	for _, word := range strings.Split(str, ",") {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		i, err := strconv.Atoi(word)
		if err == nil {
			log.SetVerbose(i)
		} else {
			log.AddKeyword(word)
		}
	}
}

func (log *DLogger) DelVerboseString(str string) {
	for _, word := range strings.Split(str, ",") {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		i, err := strconv.Atoi(word)
		if err == nil {
			log.SetVerbose(i)
		} else {
			log.DelKeyword(word)
		}
	}
}

func (log *DLogger) AddKeyword(s string) {
	if log.keywords == nil {
		log.keywords = []string{}
	}
	log.keywords = append(log.keywords, s)
}

func (log *DLogger) DelKeyword(s string) {
	for i, k := range log.keywords {
		if k == s {
			log.keywords = append(log.keywords[:i], log.keywords[i+1:]...)
		}
	}
}

func (log *DLogger) HasKeyword(s string) bool {
	for _, k := range log.keywords {
		if k == s {
			return true
		}
	}
	return false
}

// pre,post
func (l *DLogger) modIndent(grow bool) (string, string) {
	if !l.autoindent {
		return l.indent, l.indent
	}

	pre := l.indent
	l.Lock()
	defer l.Unlock()

	if grow {
		l.indent = "| " + l.indent
	} else {
		if len(l.indent) > 1 {
			l.indent = l.indent[2:]
		}
	}
	return pre, l.indent
}

func (l *DLogger) print(doit bool, f string, a ...any) {
	prefix := ""

	if len(f) > 0 && f[0] == '>' {
		if doit {
			prefix, _ = l.modIndent(true)
		}
		f = f[1:]
	} else if len(f) > 0 && f[0] == '<' {
		if doit {
			_, prefix = l.modIndent(false)
		}
		f = f[1:]
	} else {
		prefix = l.indent
	}

	if doit {
		l.log.Printf(prefix+f, a...)
	}
}

func (l *DLogger) VPrint(v int, a ...any) {
	l.print(v <= l.verbose, "%s", a...)
}

func (l *DLogger) VPrintf(v int, f string, a ...any) {
	l.print(v <= l.verbose, f, a...)
}

func (l *DLogger) VPrintln(v int, a ...any) {
	l.print(v <= l.verbose, "%s", fmt.Sprintln(a...))
}

func (l *DLogger) KPrint(key string, a ...any) {
	l.print(l.HasKeyword(key), "%s", fmt.Sprint(a...))
}

func (l *DLogger) KPrintf(key string, f string, a ...any) {
	l.print(l.HasKeyword(key), f, a...)
}

func (l *DLogger) KPrintln(key string, a ...any) {
	l.print(l.HasKeyword(key), "%s", fmt.Sprintln(a...))
}

func (l *DLogger) Trace(key any, args ...any) func() {
	if keyStr, ok := key.(string); ok {
		if !l.HasKeyword(keyStr) {
			return func() {}
		}
	} else if levelInt, ok := key.(int); ok {
		if levelInt <= l.verbose {
			return func() {}
		}
	} else if key != nil {
		panic("Trace can only be called with string or int")
	}

	fnName := ""
	doIndent := false
	saveIndent := l.indent
	var startTime time.Time
	format := ""

	if len(args) > 0 {
		ok := false
		if format, ok = args[0].(string); ok {
			pc, _, _, _ := runtime.Caller(2)
			fnName = runtime.FuncForPC(pc).Name()
			if i := strings.LastIndex(fnName, "."); i >= 0 {
				fnName = fnName[i+1:]
			}
			startTime = time.Now()

			l.modIndent(true)
			doIndent = true
			format = saveIndent + "Enter: " + fnName + ": " + format
			args = args[1:]
		}
	}

	l.log.Printf(format, args...)

	if !doIndent {
		return func() {}
	}

	return func() {
		l.modIndent(false)
		l.log.Printf(saveIndent+"Exit: %s (%v)", fnName, time.Since(startTime))
	}
}

func (l *DLogger) Fatal(a ...any)                 { l.log.Fatal(a...) }
func (l *DLogger) Fatalf(f string, a ...any)      { l.log.Fatalf(f, a...) }
func (l *DLogger) Fatalln(a ...any)               { l.log.Fatalln(a...) }
func (l *DLogger) Flags() int                     { return l.log.Flags() }
func (l *DLogger) Output(cd int, s string) error  { return l.log.Output(cd, s) }
func (l *DLogger) Panic(a ...any)                 { l.log.Panic(a...) }
func (l *DLogger) Panicf(format string, a ...any) { l.log.Panicf(format, a...) }
func (l *DLogger) Panicln(a ...any)               { l.log.Panicln(a...) }
func (l *DLogger) Prefix() string                 { return l.log.Prefix() }
func (l *DLogger) Print(a ...any)                 { l.VPrint(0, a...) }
func (l *DLogger) Printf(f string, a ...any)      { l.VPrintf(0, f, a...) }
func (l *DLogger) Println(a ...any)               { l.VPrintln(0, a...) }
func (l *DLogger) SetFlags(flag int)              { l.log.SetFlags(flag) }
func (l *DLogger) SetOutput(w io.Writer)          { l.log.SetOutput(w) }
func (l *DLogger) SetPrefix(prefix string)        { l.log.SetPrefix(prefix) }
func (l *DLogger) Writer() io.Writer              { return l.log.Writer() }

// Default logger stuff
func String() string                       { return std.String() }
func GetVerbose() int                      { return std.GetVerbose() }
func SetVerbose(v int)                     { std.SetVerbose(v) }
func SetAutoIndent(v bool)                 { std.SetAutoIndent(v) }
func AddVerboseString(str string)          { std.AddVerboseString(str) }
func DelVerboseString(str string)          { std.DelVerboseString(str) }
func AddKeyword(str string)                { std.AddKeyword(str) }
func DelKeyword(str string)                { std.DelKeyword(str) }
func HasKeyword(str string) bool           { return std.HasKeyword(str) }
func VPrint(v int, a ...any)               { std.VPrint(v, a...) }
func VPrintf(v int, f string, a ...any)    { std.VPrintf(v, f, a...) }
func VPrintln(v int, a ...any)             { std.VPrintln(v, a...) }
func KPrint(k string, a ...any)            { std.KPrint(k, a...) }
func KPrintf(k string, f string, a ...any) { std.KPrintf(k, f, a...) }
func KPrintln(k string, a ...any)          { std.KPrintln(k, a...) }
func Trace(k any, a ...any) func()         { return std.Trace(k, a...) }

func Fatal(a ...any)                 { std.Fatal(a...) }
func Fatalf(f string, a ...any)      { std.Fatalf(f, a...) }
func Fatalln(a ...any)               { std.Fatalln(a...) }
func Flags() int                     { return std.Flags() }
func Output(cd int, s string) error  { return std.Output(cd, s) }
func Panic(a ...any)                 { std.Panic(a...) }
func Panicf(format string, a ...any) { std.Panicf(format, a...) }
func Panicln(a ...any)               { std.Panicln(a...) }
func Prefix() string                 { return std.Prefix() }
func Print(a ...any)                 { std.Print(a...) }
func Printf(f string, a ...any)      { std.Printf(f, a...) }
func Println(a ...any)               { std.Println(a...) }
func SetFlags(flag int)              { std.SetFlags(flag) }
func SetOutput(w io.Writer)          { std.SetOutput(w) }
func SetPrefix(prefix string)        { std.SetPrefix(prefix) }
func Writer() io.Writer              { return std.Writer() }
