package dlog

import (
	"io"
	"log"
)

type DLogger struct {
	verbose int
	indent  string
	log     *log.Logger
}

var std = &DLogger{
	verbose: 1,
	log:     log.Default(),
}

func (log *DLogger) GetVerbose() int  { return log.verbose }
func (log *DLogger) SetVerbose(v int) { log.verbose = v }

func (l *DLogger) VPrint(v int, a ...any) {
	if v <= l.verbose {
		l.log.Print(a...)
	}
}

func (l *DLogger) VPrintf(v int, f string, a ...any) {
	saveIndent := l.indent

	if len(f) > 0 && f[0] == '>' {
		l.indent = "| " + l.indent
		f = f[1:]
		if f == "" {
			return
		}
	} else if len(f) > 0 && f[0] == '<' && len(l.indent) > 1 {
		l.indent = l.indent[2:]
		saveIndent = l.indent
		f = f[1:]
		if f == "" {
			return
		}
	}
	f = saveIndent + f

	if v <= l.verbose {
		l.log.Printf(f, a...)
	}
}

func (l *DLogger) VPrintln(v int, a ...any) {
	if v <= l.verbose {
		l.log.Println(a...)
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
func (l *DLogger) Print(a ...any)                 { l.VPrint(1, a...) }
func (l *DLogger) Printf(f string, a ...any)      { l.VPrintf(1, f, a...) }
func (l *DLogger) Println(a ...any)               { l.VPrintln(1, a...) }
func (l *DLogger) SetFlags(flag int)              { l.log.SetFlags(flag) }
func (l *DLogger) SetOutput(w io.Writer)          { l.log.SetOutput(w) }
func (l *DLogger) SetPrefix(prefix string)        { l.log.SetPrefix(prefix) }
func (l *DLogger) Writer() io.Writer              { return l.log.Writer() }

// Default logger stuff
func GetVerbose() int                   { return std.GetVerbose() }
func SetVerbose(v int)                  { std.SetVerbose(v) }
func VPrint(v int, a ...any)            { std.VPrint(v, a...) }
func VPrintf(v int, f string, a ...any) { std.VPrintf(v, f, a...) }
func VPrintln(v int, a ...any)          { std.VPrintln(v, a...) }

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
