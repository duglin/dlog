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
	log    *log.Logger
	indent string
	flat   bool // do we indent children log msgs with: |
	ascii  bool

	// Current logging flags
	currentUF *userFlags // current set of user chosen version level/triggers

	stack []*userFlags // entry created by Trace or "<"
}

var std = &DLogger{
	log: log.Default(),
	currentUF: &userFlags{
		verbose: 1,
	},
}

type userFlags struct {
	verbose       int
	showAllTraces bool
	showAllPrints bool
	showTime      bool
	triggers      map[string]*userFlags
}

func (uf *userFlags) Clone() *userFlags {
	if uf == nil {
		return nil
	}
	newUF := userFlags{
		verbose:       uf.verbose,
		showAllTraces: uf.showAllTraces,
		showAllPrints: uf.showAllPrints,
		showTime:      uf.showTime,
	}

	if len(uf.triggers) > 0 {
		newUF.triggers = map[string]*userFlags{}
		for k, kUF := range uf.triggers {
			newUF.triggers[k] = kUF.Clone()
		}
	}

	return &newUF
}

func (uf *userFlags) merge(key string, otherUF *userFlags) *userFlags {
	if otherUF.verbose > uf.verbose {
		uf.verbose = otherUF.verbose
	}
	uf.showAllTraces = uf.showAllTraces || otherUF.showAllTraces
	uf.showAllPrints = uf.showAllPrints || otherUF.showAllPrints
	uf.showTime = uf.showTime || otherUF.showTime
	for k, kw := range otherUF.triggers {
		if k == key {
			// continue
		}
		if uf.triggers == nil {
			uf.triggers = map[string]*userFlags{}
		}
		tmpkw := uf.triggers[k]
		if tmpkw != nil {
			kw = tmpkw.merge(k, kw)
		}
		uf.triggers[k] = kw
	}
	return uf
}

func (log *DLogger) Top() string {
	if log.flat {
		return ">"
	}
	if log.ascii {
		return "+"
	}
	return "┌"
}

func (log *DLogger) Bar() string {
	if log.flat {
		return ""
	}
	if log.ascii {
		return "|"
	}
	return "│"
}

func (log *DLogger) Bot() string {
	if log.flat {
		return "<"
	}
	if log.ascii {
		return "+"
	}
	return "└"
}

const TOP = "┌"

func (log *DLogger) String() string {
	str := ""
	if !log.flat {
		str += "Indented,"
	}
	if log.ascii {
		str += "Ascii,"
	}

	str += "Current: "

	fn := func(string, string, *userFlags) string { return "" }
	fn = func(indent string, key string, uf *userFlags) string {
		s := fmt.Sprintf("%d", uf.verbose)
		if uf.showAllTraces {
			s += ", AllTraces"
		}
		if uf.showAllPrints {
			s += ", AllPrints"
		}
		if uf.showTime {
			s += ", ShowTimes"
		}
		s += "\n"

		indent = strings.Repeat(" ", len(indent))
		for key, nextUF := range uf.triggers {
			s += indent + "  L " + key + ": " + fn(indent+"  ", key, nextUF)
		}
		return s
	}

	return str + fn("Current:", "", log.currentUF)
}

func (log *DLogger) Dump() string {
	buf := &strings.Builder{}

	if !log.flat {
		fmt.Fprint(buf, "indent,")
	}

	if log.ascii {
		fmt.Fprint(buf, "ascii,")
	}

	fmt.Fprintf(buf, "curr(")

	fn := func(*userFlags) {}
	fn = func(uf *userFlags) {
		if uf == nil {
			fmt.Fprintf(buf, "nil")
			return
		}

		fmt.Fprintf(buf, "v:%d", uf.verbose)
		if uf.showAllTraces {
			fmt.Fprintf(buf, ",traces")
		}
		if uf.showAllPrints {
			fmt.Fprintf(buf, ",prints")
		}
		if uf.showTime {
			fmt.Fprintf(buf, ",time")
		}

		i := 0
		for k, kUF := range uf.triggers {
			if i == 0 {
				fmt.Fprintf(buf, ",[")
			}
			fmt.Fprintf(buf, "%s(", k)
			fn(kUF)
			fmt.Fprintf(buf, ")")
			if i+1 == len(uf.triggers) {
				fmt.Fprintf(buf, "]")
			} else {
				fmt.Fprintf(buf, ",")
			}
			i++
		}
	}

	fn(log.currentUF)
	fmt.Fprintf(buf, ")")
	/*
		        if len(log.stack) > 0 {
				    fmt.Fprintf(buf, ",stack(")
				for _, s := range log.stack {
					fn(s)
				}
		            fmt.Fprintf(buf,")")
		        }
	*/
	return buf.String()
}

func (log *DLogger) Reset() {
	log.indent = ""
	log.flat = false
	log.ascii = false
	log.currentUF = &userFlags{}
	log.stack = nil
}

func (log *DLogger) GetVerbose() int  { return log.currentUF.verbose }
func (log *DLogger) SetIndent(v bool) { log.flat = !v }
func (log *DLogger) SetAscii(v bool)  { log.ascii = v }

// trigger must be: int or string (xxx,xxx where xxx is int or string)
func (log *DLogger) SetVerbose(triggers ...any) {
	for _, trigger := range triggers {
		v, ok := trigger.(int)
		if ok {
			log.currentUF.verbose = v
			continue
		}
		str, ok := trigger.(string)
		if !ok {
			panic(fmt.Sprintf("trigger must be an int or string: %v", trigger))
		}

		for _, expr := range strings.Split(str, ",") {
			if expr = strings.TrimSpace(expr); expr == "" {
				continue
			}

			fn := func(*userFlags, []string) {}
			fn = func(uf *userFlags, words []string) {
				for i, word := range words {
					word = strings.TrimSpace(word)
					if word == "" {
						continue
					} else if word == "flat" {
						log.flat = true
					} else if word == "ascii" {
						log.ascii = true
					} else if word == "time" || word == "timed" {
						uf.showTime = true
					} else if word == "*" {
						uf.showAllTraces = true
					} else if word == "**" {
						uf.showAllPrints = true
					} else if v, err := strconv.Atoi(word); err == nil {
						uf.verbose = v
					} else {
						if uf.triggers == nil {
							uf.triggers = map[string]*userFlags{}
						}
						newUF := uf.triggers[word]
						if newUF == nil {
							newUF = &userFlags{}
							uf.triggers[word] = newUF
						}
						fn(newUF, words[i+1:])
						break
					}
				}
			}

			words := strings.Split(expr, ":")
			fn(log.currentUF, words)
		}
	}
}

func (log *DLogger) DelVerbose(triggers ...any) {
	rootUF := log.currentUF
	if l := len(log.stack); l > 0 {
		rootUF = log.stack[l-1]
	}

	for _, trigger := range triggers {
		if _, ok := trigger.(int); ok {
			log.currentUF.verbose = 0
			continue
		}

		str, ok := trigger.(string)
		if !ok {
			panic(fmt.Sprintf("trigger must be an int or string: %v", trigger))
		}

		for _, expr := range strings.Split(str, ",") {
			words := strings.Split(expr, ":")
			tmpUF := rootUF
			numWords := len(words)
			for i, word := range words {
				if word = strings.TrimSpace(word); word == "" {
					continue
				}

				if word == "*" {
					tmpUF.showAllTraces = false
				} else if word == "**" {
					tmpUF.showAllPrints = false
				} else {
					if i+1 == numWords {
						if _, intErr := strconv.Atoi(word); intErr == nil {
							tmpUF.verbose = 0
						} else {
							delete(tmpUF.triggers, word)
						}
					} else {
						tmpUF = tmpUF.triggers[word]
						if tmpUF == nil {
							break // stop traversing
						}
					}
				}
			}
		}
	}
}

// is "s" enabled as a verbose keyword?
func (log *DLogger) HasVerbose(s string) bool {
	return log.currentUF.triggers[s] != nil
}

var pc2fn = map[uintptr]string{}

// Is this function's name in the list of verbose keywords?
func (log *DLogger) isFuncVerboseDepth(d int) bool {
	fnName := ""

	pc, _, _, _ := runtime.Caller(d)
	if fnName, _ = pc2fn[pc]; fnName == "" {
		// Not seen this PC yet so get its fnName
		fnName = runtime.FuncForPC(pc).Name()
		if i := strings.LastIndex(fnName, "."); i >= 0 {
			fnName = fnName[i+1:]
		}
		pc2fn[pc] = fnName
	}

	return log.currentUF.triggers[fnName] != nil
}

// Is this function's name in the list of verbose keywords?
func (log *DLogger) IsFuncVerbose() bool {
	if log.currentUF.showAllPrints {
		return true
	}

	if len(log.currentUF.triggers) == 0 {
		return false
	}

	return log.isFuncVerboseDepth(3)
}

// Is ANY keyword verbose trigger set?
func (log *DLogger) IsVerbose() bool {
	return len(log.currentUF.triggers) > 0
}

// pre,post indent string
func (l *DLogger) modScope(grow bool, key string) (string, string) {
	pre := l.indent

	if !l.flat {
		l.Lock()
		defer l.Unlock()
	}

	indentStr := l.Bar() + "   "
	if grow {
		if !l.flat {
			l.indent = indentStr + l.indent
		}
		// fmt.Printf("PRE MOD (%s):\n%s\n", key, l.Dump())

		l.stack = append([]*userFlags{l.currentUF}, l.stack...)
		l.currentUF = l.currentUF.triggers[key].Clone()
		// fmt.Printf("POST-1 MOD (%s):\n%s\n", key, l.Dump())
		if l.currentUF == nil {
			l.currentUF = l.stack[0].Clone()
		} else {
			l.currentUF.merge(key, l.stack[0])
		}
		// fmt.Printf("POST-2 MOD (%s):\n%s\n", key, l.Dump())
	} else {
		if !l.flat && len(l.indent) > 1 {
			l.indent = l.indent[len(indentStr):]
		}
		if len(l.stack) > 0 {
			l.currentUF = l.stack[0]
			l.stack = l.stack[1:]
		}
	}

	return pre, l.indent
}

func (l *DLogger) checkVTrigger(key any) bool {
	if keyStr, ok := key.(string); ok {
		return l.HasVerbose(keyStr)
	} else if levelInt, ok := key.(int); ok {
		return levelInt <= l.currentUF.verbose
	} else if key != nil {
		return true
	}

	panic("Verbose trigger can only be called with string or int")
}

func (l *DLogger) VPrint(k any, a ...any) func() {
	if l.currentUF.showAllPrints || l.checkVTrigger(k) {
		l.log.Printf(l.indent + JoinArgs(a...))
	}
	return func() {}
}

func (l *DLogger) VPrintf(k any, f string, a ...any) func() {
	if l.currentUF.showAllPrints || l.checkVTrigger(k) {
		l.log.Printf(l.indent+f, a...)
	}
	return func() {}
}

func (l *DLogger) VPrintln(k any, a ...any) func() {
	if l.currentUF.showAllPrints || l.checkVTrigger(k) {
		l.log.Printf(l.indent + fmt.Sprintln(a...))
	}
	return func() {}

}

func (l *DLogger) FuncPrint(a ...any) func() {
	if !l.currentUF.showAllPrints && len(l.currentUF.triggers) == 0 {
		return func() {}
	}

	if !l.isFuncVerboseDepth(3) {
		return func() {}
	}

	l.log.Printf(l.indent + JoinArgs(a...))
	return func() {}
}

func (l *DLogger) FuncPrintf(f string, a ...any) func() {
	if !l.currentUF.showAllPrints && len(l.currentUF.triggers) == 0 {
		return func() {}
	}

	if !l.isFuncVerboseDepth(3) {
		return func() {}
	}

	l.log.Printf(l.indent+f, a...)
	return func() {}
}

func (l *DLogger) FuncPrintln(a ...any) func() {
	if !l.currentUF.showAllPrints && len(l.currentUF.triggers) == 0 {
		return func() {}
	}

	if !l.isFuncVerboseDepth(3) {
		return func() {}
	}

	l.log.Printf(l.indent + fmt.Sprintln(a...))
	return func() {}
}

func (l *DLogger) VTracef(key any, f string, args ...any) func() {
	return VTrace(key, append([]any{f}, args...)...)
}

func (l *DLogger) VTrace(key any, args ...any) func() {
	if !l.currentUF.showAllTraces && !l.checkVTrigger(key) {
		return func() {}
	}

	modKey, ok := key.(string)
	if !ok {
		modKey = fmt.Sprintf("%v", key)
	}
	showKey := modKey

	format := l.indent
	ENTER := l.Top() + " " // "+ " // "+ " // ENTER := "Enter: "
	EXIT := l.Bot() + " "  // "+ "  // "+ "  // "└─ " // EXIT := "Exit: "

	format += ENTER + showKey // + ":"
	if len(showKey) > 0 && len(args) > 0 {
		format += ":"
	}
	startTime := time.Now()

	l.modScope(true, modKey)
	doTime := l.currentUF.showTime

	retFunc := func() {
		_, indent := l.modScope(false, "")
		suffix := ""
		if doTime {
			suffix = fmt.Sprintf(" time: %v", time.Since(startTime))
		}
		l.log.Printf(indent + EXIT + showKey + suffix)
	}

	if len(args) > 0 {
		if fmtStr, ok := args[0].(string); ok {
			if strings.Contains(fmtStr, "%") {
				format += " " + fmt.Sprintf(fmtStr, args[1:]...)
			} else {
				format += " " + JoinArgs(args...)
			}
		} else {
			format += " " + JoinArgs(args...)
		}
	}

	l.log.Printf(format)
	return retFunc
}

func (l *DLogger) Tracef(f string, args ...any) func() {
	// return VTrace(append([]any{f}, args...)...)
	if !l.currentUF.showAllTraces && len(l.currentUF.triggers) == 0 {
		return func() {}
	}

	fnName := ""
	pc, _, _, _ := runtime.Caller(2)
	if fnName, _ = pc2fn[pc]; fnName == "" {
		fnName = runtime.FuncForPC(pc).Name()
		if i := strings.LastIndex(fnName, "."); i >= 0 {
			fnName = fnName[i+1:]
		}
		pc2fn[pc] = fnName
	}
	return VTrace(fnName, append([]any{f}, args...)...)
}

func (l *DLogger) Trace(args ...any) func() {
	if !l.currentUF.showAllTraces && len(l.currentUF.triggers) == 0 {
		return func() {}
	}

	fnName := ""
	pc, _, _, _ := runtime.Caller(2)
	if fnName, _ = pc2fn[pc]; fnName == "" {
		fnName = runtime.FuncForPC(pc).Name()
		if i := strings.LastIndex(fnName, "."); i >= 0 {
			fnName = fnName[i+1:]
		}
		pc2fn[pc] = fnName
	}

	return VTrace(fnName, args...)
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
func Reset()                   { std.Reset() }
func GetVerbose() int          { return std.GetVerbose() }
func SetVerbose(k ...any)      { std.SetVerbose(k...) }
func DelVerbose(k ...any)      { std.DelVerbose(k...) }
func HasVerbose(k string) bool { return std.HasVerbose(k) }
func IsFuncVerbose() bool      { return std.IsFuncVerbose() }
func IsVerbose() bool          { return std.IsVerbose() }
func SetIndent(v bool)         { std.SetIndent(v) }
func String() string           { return std.String() }
func Dump() string             { return std.Dump() }

func VPrint(k any, a ...any) func()            { return std.VPrint(k, a...) }
func VPrintf(k any, f string, a ...any) func() { return std.VPrintf(k, f, a...) }
func VPrintln(k any, a ...any) func()          { return std.VPrintln(k, a...) }
func Trace(a ...any) func()                    { return std.Trace(a...) }
func Tracef(f string, a ...any) func()         { return std.Tracef(f, a...) }
func VTrace(k any, a ...any) func()            { return std.VTrace(k, a...) }
func VTracef(k any, f string, a ...any) func() { return std.VTracef(k, f, a...) }

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
func FuncPrint(a ...any)             { std.FuncPrint(a...) }
func FuncPrintf(f string, a ...any)  { std.FuncPrintf(f, a...) }
func FuncPrintln(a ...any)           { std.FuncPrintln(a...) }
func SetFlags(flag int)              { std.SetFlags(flag) }
func SetOutput(w io.Writer)          { std.SetOutput(w) }
func SetPrefix(prefix string)        { std.SetPrefix(prefix) }
func Writer() io.Writer              { return std.Writer() }

func JoinArgs(args ...any) string {
	buf := &strings.Builder{}
	for _, a := range args {
		if buf.Len() != 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(buf, "%v", a)
	}
	return buf.String()
}
