package dlog

import (
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
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

func (l *DLogger) VPrint(v int, f string, a ...any) {
	f = l.processIndents(f)
	if f == "" || v > l.verbose {
		return
	}
	l.log.Print(a...)
}

func (l *DLogger) processIndents(f string) string {
	if len(f) == 0 {
		return ""
	}

	l.Lock()
	saveIndent := l.indent
	l.Unlock()
	iLen := len(saveIndent)

	if len(f) > 0 && f[0] == '>' {
		if l.autoindent {
			l.Lock()
			l.indent = "| " + l.indent
			l.Unlock()
		}
		f = f[1:]
		if f == "" {
			return f
		}
	} else if len(f) > 0 && f[0] == '<' && iLen > 1 {
		if l.autoindent {
			l.Lock()
			l.indent = l.indent[2:]
			saveIndent = l.indent
			l.Unlock()
		}
		f = f[1:]
		if f == "" {
			return f
		}
	}

	return saveIndent + f
}

func (l *DLogger) VPrintf(v int, f string, a ...any) {
	f = l.processIndents(f)
	if f == "" || v > l.verbose {
		return
	}

	l.log.Printf(f, a...)
}

func (l *DLogger) VPrintln(v int, f string, a ...any) {
	f = l.processIndents(f)
	if f == "" || v > l.verbose {
		return
	}

	l.log.Println(a...)
}

func (l *DLogger) KPrint(key string, f string, a ...any) {
	f = l.processIndents(f)
	if f == "" || !l.HasKeyword(key) {
		return
	}

	a = append([]any{f}, a...)
	l.log.Print(a...)
}

func (l *DLogger) KPrintf(key string, f string, a ...any) {
	f = l.processIndents(f)
	if f == "" || !l.HasKeyword(key) {
		return
	}
	l.log.Printf(f, a...)
}

func (l *DLogger) KPrintln(key string, f string, a ...any) {
	f = l.processIndents(f)
	if f == "" || !l.HasKeyword(key) {
		return
	}
	a = append([]any{f}, a...)
	l.log.Println(a...)
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
func (l *DLogger) Print(f string, a ...any)       { l.VPrint(0, f, a...) }
func (l *DLogger) Printf(f string, a ...any)      { l.VPrintf(0, f, a...) }
func (l *DLogger) Println(f string, a ...any)     { l.VPrintln(0, f, a...) }
func (l *DLogger) SetFlags(flag int)              { l.log.SetFlags(flag) }
func (l *DLogger) SetOutput(w io.Writer)          { l.log.SetOutput(w) }
func (l *DLogger) SetPrefix(prefix string)        { l.log.SetPrefix(prefix) }
func (l *DLogger) Writer() io.Writer              { return l.log.Writer() }

// Default logger stuff
func GetVerbose() int                       { return std.GetVerbose() }
func SetVerbose(v int)                      { std.SetVerbose(v) }
func SetAutoIndent(v bool)                  { std.SetAutoIndent(v) }
func HasKeyword(str string) bool            { return std.HasKeyword(str) }
func AddVerboseString(str string)           { std.AddVerboseString(str) }
func DelVerboseString(str string)           { std.DelVerboseString(str) }
func VPrint(v int, f string, a ...any)      { std.VPrint(v, f, a...) }
func VPrintf(v int, f string, a ...any)     { std.VPrintf(v, f, a...) }
func VPrintln(v int, f string, a ...any)    { std.VPrintln(v, f, a...) }
func KPrint(k string, f string, a ...any)   { std.KPrint(k, f, a...) }
func KPrintf(k string, f string, a ...any)  { std.KPrintf(k, f, a...) }
func KPrintln(k string, f string, a ...any) { std.KPrintln(k, f, a...) }

func Fatal(a ...any)                 { std.Fatal(a...) }
func Fatalf(f string, a ...any)      { std.Fatalf(f, a...) }
func Fatalln(a ...any)               { std.Fatalln(a...) }
func Flags() int                     { return std.Flags() }
func Output(cd int, s string) error  { return std.Output(cd, s) }
func Panic(a ...any)                 { std.Panic(a...) }
func Panicf(format string, a ...any) { std.Panicf(format, a...) }
func Panicln(a ...any)               { std.Panicln(a...) }
func Prefix() string                 { return std.Prefix() }
func Print(f string, a ...any)       { std.Print(f, a...) }
func Printf(f string, a ...any)      { std.Printf(f, a...) }
func Println(f string, a ...any)     { std.Println(f, a...) }
func SetFlags(flag int)              { std.SetFlags(flag) }
func SetOutput(w io.Writer)          { std.SetOutput(w) }
func SetPrefix(prefix string)        { std.SetPrefix(prefix) }
func Writer() io.Writer              { return std.Writer() }
