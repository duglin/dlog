package test

import (
	// "fmt"
	"regexp"
	"strings"
	"testing"

	log "github.com/duglin/dlog"
)

func doit(str string, re bool) {
	defer log.Trace("str: %q, recursive: %v", str, re)()
	log.Trace("no defer", str, re)()
	fn := log.Trace("fn", str, re)

	log.Printf("printf from doit: %s / %v", str, re)
	log.VPrintf(0, "0-vprintf from doit: %s / %v", str, re)
	log.VPrintf(2, "2-vprintf from doit: %s / %v", str, re)
	log.VPrintf("showme", "showme-vprintf from doit: %s / %v", str, re)

	fn()

	if re {
		doit(str+"-doit", !re)
	}
}

func run() {
	defer log.Trace()()

	log.Printf("printf from run")
	log.VPrintf(0, "0-vprintf from run")
	log.VPrintf(2, "2-vprintf from run")
	log.VPrintf("showme", "showme-vprintf from run")
	doit("run", true)
}

func main() {
	defer log.Trace()()
	log.Printf("printf from main")
	log.VPrintf(0, "0-vprintf from main")
	log.VPrintf(2, "2-vprintf from main")
	log.VPrintf("showme", "showme-vprintf from main")
	run()
	doit("main", false)
	log.Printf("exiting from main")
}

func runScoped(from string) {
	defer log.Tracef("from: %s", from)()
	defer log.Tracef(from)()
	fn := log.VTracef("rr", "rr-v from runScoped - from: %s", from)
	log.VPrintf(1, "1-vprintf from runScoped - from: %s", from)
	fn()
	log.VPrintf(2, "2-vprintf from runScoped - from: %s", from)
}

func doitScoped() {
	defer log.Trace()()
	log.VPrintf(1, "1-vprintf from doitScoped")
	log.VPrintf("show", "show from doitScoped")
	runScoped("doitScoped")
}

func mainScoped() {
	defer log.Trace()()
	fn := log.VTrace(1, "1-vprintf from mainScoped - enter")
	log.VPrint(1, "1-vprintf from mainScoped")
	log.VPrintf("mm", "mm-vprintf from mainScoped")
	doitScoped()
	fn()
	fn = log.Trace()
	runScoped("mainScoped")
	fn()
}

func mainVerboseFunc() {
	if log.VerboseFunc() {
		log.Printf("I'm only shown when mainVerboseFunc is on")
	}
}

var REG_TS = `\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `
var TIME_TS = `time: \d+\.\d+.?s`

func diff(t *testing.T, got *strings.Builder, exp string) {
	t.Helper()
	// Grab output and mask it
	gotStr := got.String()
	rawStr := gotStr
	gotStr = regexp.MustCompile("(?s)"+REG_TS).ReplaceAllString(gotStr, "TS ")
	gotStr = regexp.MustCompile("(?s)"+TIME_TS).ReplaceAllString(gotStr, "time: ss")

	exp = regexp.MustCompile("(?s)"+REG_TS).ReplaceAllString(exp, "TS ")
	exp = regexp.MustCompile("(?s)"+TIME_TS).ReplaceAllString(exp, "time: ss")
	// time: 14.917µs

	if exp != "" && exp[0] == '\n' {
		exp = exp[1:]
	}

	res := Diff("Exp", exp, "Got", gotStr)

	if res != "" {
		t.Fatalf("\nGOT:\n%s\nEXP:\n%s\n%s", rawStr, exp, res)
	}

	got.Reset()
}
