package test

import (
	// "fmt"
	"strings"
	"testing"

	log "github.com/duglin/dlog"
	golog "log"
)

func TestBasic(t *testing.T) {
	res := &strings.Builder{}
	golog.SetOutput(res)

	// No logging outout
	log.SetIndent(false)
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS printf from run
TS 0-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS exiting from main
`)

	// Set verbose=1, no change since we have no logginf for level 1
	log.SetVerbose(1)
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS printf from run
TS 0-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS exiting from main
`)

	// Now set verbose-2, we do have some logging now
	log.SetVerbose(2)
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS 2-vprintf from main
TS printf from run
TS 0-vprintf from run
TS 2-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS 2-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS 2-vprintf from doit: run-doit / false
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS 2-vprintf from doit: main / false
TS exiting from main
`)

	// Stop 1+ logging, but show log for "main" entry/exit
	log.SetVerbose(0)
	log.SetVerbose("main")
	main()
	diff(t, res, `
TS > main
TS printf from main
TS 0-vprintf from main
TS printf from run
TS 0-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS exiting from main
TS < main
`)

	// Reset + trace run()
	log.Reset()
	log.SetIndent(false)
	log.SetVerbose("run")
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS > run
TS printf from run
TS 0-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS < run
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS exiting from main
`)

	// Add indent + main
	log.Reset()
	log.SetVerbose("main")
	main()
	diff(t, res, `
TS ┌ main
TS │   printf from main
TS │   0-vprintf from main
TS │   printf from run
TS │   0-vprintf from run
TS │   printf from doit: run / true
TS │   0-vprintf from doit: run / true
TS │   printf from doit: run-doit / false
TS │   0-vprintf from doit: run-doit / false
TS │   printf from doit: main / false
TS │   0-vprintf from doit: main / false
TS │   exiting from main
TS └ main
`)

	// Add main+run
	log.Reset()
	log.SetVerbose("main:run")
	main()
	diff(t, res, `
TS ┌ main
TS │   printf from main
TS │   0-vprintf from main
TS │   ┌ run
TS │   │   printf from run
TS │   │   0-vprintf from run
TS │   │   printf from doit: run / true
TS │   │   0-vprintf from doit: run / true
TS │   │   printf from doit: run-doit / false
TS │   │   0-vprintf from doit: run-doit / false
TS │   └ run
TS │   printf from doit: main / false
TS │   0-vprintf from doit: main / false
TS │   exiting from main
TS └ main
`)

	// Add main:doit (no run)
	log.Reset()
	log.SetVerbose("main:doit")
	main()
	diff(t, res, `
TS ┌ main
TS │   printf from main
TS │   0-vprintf from main
TS │   printf from run
TS │   0-vprintf from run
TS │   ┌ doit: str: "run", recursive: true
TS │   │   ┌ doit: no defer run true
TS │   │   └ doit
TS │   │   ┌ doit: fn run true
TS │   │   │   printf from doit: run / true
TS │   │   │   0-vprintf from doit: run / true
TS │   │   └ doit
TS │   │   ┌ doit: str: "run-doit", recursive: false
TS │   │   │   ┌ doit: no defer run-doit false
TS │   │   │   └ doit
TS │   │   │   ┌ doit: fn run-doit false
TS │   │   │   │   printf from doit: run-doit / false
TS │   │   │   │   0-vprintf from doit: run-doit / false
TS │   │   │   └ doit
TS │   │   └ doit
TS │   └ doit
TS │   ┌ doit: str: "main", recursive: false
TS │   │   ┌ doit: no defer main false
TS │   │   └ doit
TS │   │   ┌ doit: fn main false
TS │   │   │   printf from doit: main / false
TS │   │   │   0-vprintf from doit: main / false
TS │   │   └ doit
TS │   └ doit
TS │   exiting from main
TS └ main
`)

	// Add main,doit (no run), same as before
	log.Reset()
	log.SetVerbose("main,doit")
	main()
	diff(t, res, `
TS ┌ main
TS │   printf from main
TS │   0-vprintf from main
TS │   printf from run
TS │   0-vprintf from run
TS │   ┌ doit: str: "run", recursive: true
TS │   │   ┌ doit: no defer run true
TS │   │   └ doit
TS │   │   ┌ doit: fn run true
TS │   │   │   printf from doit: run / true
TS │   │   │   0-vprintf from doit: run / true
TS │   │   └ doit
TS │   │   ┌ doit: str: "run-doit", recursive: false
TS │   │   │   ┌ doit: no defer run-doit false
TS │   │   │   └ doit
TS │   │   │   ┌ doit: fn run-doit false
TS │   │   │   │   printf from doit: run-doit / false
TS │   │   │   │   0-vprintf from doit: run-doit / false
TS │   │   │   └ doit
TS │   │   └ doit
TS │   └ doit
TS │   ┌ doit: str: "main", recursive: false
TS │   │   ┌ doit: no defer main false
TS │   │   └ doit
TS │   │   ┌ doit: fn main false
TS │   │   │   printf from doit: main / false
TS │   │   │   0-vprintf from doit: main / false
TS │   │   └ doit
TS │   └ doit
TS │   exiting from main
TS └ main
`)

	// run:doit (no main:doit)
	log.Reset()
	log.SetVerbose("run:doit")
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS ┌ run
TS │   printf from run
TS │   0-vprintf from run
TS │   ┌ doit: str: "run", recursive: true
TS │   │   ┌ doit: no defer run true
TS │   │   └ doit
TS │   │   ┌ doit: fn run true
TS │   │   │   printf from doit: run / true
TS │   │   │   0-vprintf from doit: run / true
TS │   │   └ doit
TS │   │   ┌ doit: str: "run-doit", recursive: false
TS │   │   │   ┌ doit: no defer run-doit false
TS │   │   │   └ doit
TS │   │   │   ┌ doit: fn run-doit false
TS │   │   │   │   printf from doit: run-doit / false
TS │   │   │   │   0-vprintf from doit: run-doit / false
TS │   │   │   └ doit
TS │   │   └ doit
TS │   └ doit
TS └ run
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS exiting from main
`)

	// doit:time
	log.Reset()
	log.SetVerbose("doit:time")
	main()
	diff(t, res, `
2026/08/02 18:59:51 printf from main
2026/08/02 18:59:51 0-vprintf from main
2026/08/02 18:59:51 printf from run
2026/08/02 18:59:51 0-vprintf from run
2026/08/02 18:59:51 ┌ doit: str: "run", recursive: true
2026/08/02 18:59:51 │   ┌ doit: no defer run true
2026/08/02 18:59:51 │   └ doit time: 1.417µs
2026/08/02 18:59:51 │   ┌ doit: fn run true
2026/08/02 18:59:51 │   │   printf from doit: run / true
2026/08/02 18:59:51 │   │   0-vprintf from doit: run / true
2026/08/02 18:59:51 │   └ doit time: 9.731µs
2026/08/02 18:59:51 │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 18:59:51 │   │   ┌ doit: no defer run-doit false
2026/08/02 18:59:51 │   │   └ doit time: 1.306µs
2026/08/02 18:59:51 │   │   ┌ doit: fn run-doit false
2026/08/02 18:59:51 │   │   │   printf from doit: run-doit / false
2026/08/02 18:59:51 │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 18:59:51 │   │   └ doit time: 1.804µs
2026/08/02 18:59:51 │   └ doit time: 11.067µs
2026/08/02 18:59:51 └ doit time: 58.569µs
2026/08/02 18:59:51 ┌ doit: str: "main", recursive: false
2026/08/02 18:59:51 │   ┌ doit: no defer main false
2026/08/02 18:59:51 │   └ doit time: 1.187µs
2026/08/02 18:59:51 │   ┌ doit: fn main false
2026/08/02 18:59:51 │   │   printf from doit: main / false
2026/08/02 18:59:51 │   │   0-vprintf from doit: main / false
2026/08/02 18:59:51 │   └ doit time: 4.93µs
2026/08/02 18:59:51 └ doit time: 19.523µs
2026/08/02 18:59:51 exiting from main
`)

	// run:time:* - time for run&doit, but not main's doit
	log.Reset()
	log.SetVerbose("run:time:*")
	main()
	diff(t, res, `
2026/08/02 19:00:59 printf from main
2026/08/02 19:00:59 0-vprintf from main
2026/08/02 19:00:59 ┌ run
2026/08/02 19:00:59 │   printf from run
2026/08/02 19:00:59 │   0-vprintf from run
2026/08/02 19:00:59 │   ┌ doit: str: "run", recursive: true
2026/08/02 19:00:59 │   │   ┌ doit: no defer run true
2026/08/02 19:00:59 │   │   └ doit time: 1.507µs
2026/08/02 19:00:59 │   │   ┌ doit: fn run true
2026/08/02 19:00:59 │   │   │   printf from doit: run / true
2026/08/02 19:00:59 │   │   │   0-vprintf from doit: run / true
2026/08/02 19:00:59 │   │   └ doit time: 6.206µs
2026/08/02 19:00:59 │   │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 19:00:59 │   │   │   ┌ doit: no defer run-doit false
2026/08/02 19:00:59 │   │   │   └ doit time: 5.929µs
2026/08/02 19:00:59 │   │   │   ┌ doit: fn run-doit false
2026/08/02 19:00:59 │   │   │   │   printf from doit: run-doit / false
2026/08/02 19:00:59 │   │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 19:00:59 │   │   │   └ doit time: 2.002µs
2026/08/02 19:00:59 │   │   └ doit time: 12.648µs
2026/08/02 19:00:59 │   └ doit time: 31.298µs
2026/08/02 19:00:59 └ run time: 34.765µs
2026/08/02 19:00:59 printf from doit: main / false
2026/08/02 19:00:59 0-vprintf from doit: main / false
2026/08/02 19:00:59 exiting from main
`)

	// run:doit:time,main:doit:time - so everyone except run() does time
	log.Reset()
	log.SetVerbose("run:doit:time,main:doit:time")
	main()
	diff(t, res, `
2026/08/02 19:01:31 ┌ main
2026/08/02 19:01:31 │   printf from main
2026/08/02 19:01:31 │   0-vprintf from main
2026/08/02 19:01:31 │   ┌ run
2026/08/02 19:01:31 │   │   printf from run
2026/08/02 19:01:31 │   │   0-vprintf from run
2026/08/02 19:01:31 │   │   ┌ doit: str: "run", recursive: true
2026/08/02 19:01:31 │   │   │   ┌ doit: no defer run true
2026/08/02 19:01:31 │   │   │   └ doit time: 1.37µs
2026/08/02 19:01:31 │   │   │   ┌ doit: fn run true
2026/08/02 19:01:31 │   │   │   │   printf from doit: run / true
2026/08/02 19:01:31 │   │   │   │   0-vprintf from doit: run / true
2026/08/02 19:01:31 │   │   │   └ doit time: 2.093µs
2026/08/02 19:01:31 │   │   │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 19:01:31 │   │   │   │   ┌ doit: no defer run-doit false
2026/08/02 19:01:31 │   │   │   │   └ doit time: 8.054µs
2026/08/02 19:01:31 │   │   │   │   ┌ doit: fn run-doit false
2026/08/02 19:01:31 │   │   │   │   │   printf from doit: run-doit / false
2026/08/02 19:01:31 │   │   │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 19:01:31 │   │   │   │   └ doit time: 5.177µs
2026/08/02 19:01:31 │   │   │   └ doit time: 17.762µs
2026/08/02 19:01:31 │   │   └ doit time: 37.719µs
2026/08/02 19:01:31 │   └ run
2026/08/02 19:01:31 │   ┌ doit: str: "main", recursive: false
2026/08/02 19:01:31 │   │   ┌ doit: no defer main false
2026/08/02 19:01:31 │   │   └ doit time: 4.318µs
2026/08/02 19:01:31 │   │   ┌ doit: fn main false
2026/08/02 19:01:31 │   │   │   printf from doit: main / false
2026/08/02 19:01:31 │   │   │   0-vprintf from doit: main / false
2026/08/02 19:01:31 │   │   └ doit time: 1.915µs
2026/08/02 19:01:31 │   └ doit time: 10.568µs
2026/08/02 19:01:31 │   exiting from main
2026/08/02 19:01:31 └ main
`)

	// show all logging (all levels) but just for run:doit, not main:doit
	log.Reset()
	log.SetVerbose("run:doit:**")
	main()
	diff(t, res, `
2026/08/02 21:05:50 printf from main
2026/08/02 21:05:50 0-vprintf from main
2026/08/02 21:05:50 ┌ run
2026/08/02 21:05:50 │   printf from run
2026/08/02 21:05:50 │   0-vprintf from run
2026/08/02 21:05:50 │   ┌ doit: str: "run", recursive: true
2026/08/02 21:05:50 │   │   ┌ doit: no defer run true
2026/08/02 21:05:50 │   │   └ doit
2026/08/02 21:05:50 │   │   ┌ doit: fn run true
2026/08/02 21:05:50 │   │   │   printf from doit: run / true
2026/08/02 21:05:50 │   │   │   0-vprintf from doit: run / true
2026/08/02 21:05:50 │   │   │   2-vprintf from doit: run / true
2026/08/02 21:05:50 │   │   │   showme-vprintf from doit: run / true
2026/08/02 21:05:50 │   │   └ doit
2026/08/02 21:05:50 │   │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 21:05:50 │   │   │   ┌ doit: no defer run-doit false
2026/08/02 21:05:50 │   │   │   └ doit
2026/08/02 21:05:50 │   │   │   ┌ doit: fn run-doit false
2026/08/02 21:05:50 │   │   │   │   printf from doit: run-doit / false
2026/08/02 21:05:50 │   │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 21:05:50 │   │   │   │   2-vprintf from doit: run-doit / false
2026/08/02 21:05:50 │   │   │   │   showme-vprintf from doit: run-doit / false
2026/08/02 21:05:50 │   │   │   └ doit
2026/08/02 21:05:50 │   │   └ doit
2026/08/02 21:05:50 │   └ doit
2026/08/02 21:05:50 └ run
2026/08/02 21:05:50 printf from doit: main / false
2026/08/02 21:05:50 0-vprintf from doit: main / false
2026/08/02 21:05:50 exiting from main
`)

	// all logs for run:doit, merged with main:doit:timed
	// notice not all logging when main->doit is called
	log.Reset()
	log.SetVerbose("run:doit:**,main:doit:timed")
	// fmt.Printf("TEST: %s\n", log.Dump())
	main()
	diff(t, res, `
2026/08/02 21:06:33 ┌ main
2026/08/02 21:06:33 │   printf from main
2026/08/02 21:06:33 │   0-vprintf from main
2026/08/02 21:06:33 │   ┌ run
2026/08/02 21:06:33 │   │   printf from run
2026/08/02 21:06:33 │   │   0-vprintf from run
2026/08/02 21:06:33 │   │   ┌ doit: str: "run", recursive: true
2026/08/02 21:06:33 │   │   │   ┌ doit: no defer run true
2026/08/02 21:06:33 │   │   │   └ doit time: 1.764µs
2026/08/02 21:06:33 │   │   │   ┌ doit: fn run true
2026/08/02 21:06:33 │   │   │   │   printf from doit: run / true
2026/08/02 21:06:33 │   │   │   │   0-vprintf from doit: run / true
2026/08/02 21:06:33 │   │   │   │   2-vprintf from doit: run / true
2026/08/02 21:06:33 │   │   │   │   showme-vprintf from doit: run / true
2026/08/02 21:06:33 │   │   │   └ doit time: 3.024µs
2026/08/02 21:06:33 │   │   │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 21:06:33 │   │   │   │   ┌ doit: no defer run-doit false
2026/08/02 21:06:33 │   │   │   │   └ doit time: 5.03µs
2026/08/02 21:06:33 │   │   │   │   ┌ doit: fn run-doit false
2026/08/02 21:06:33 │   │   │   │   │   printf from doit: run-doit / false
2026/08/02 21:06:33 │   │   │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 21:06:33 │   │   │   │   │   2-vprintf from doit: run-doit / false
2026/08/02 21:06:33 │   │   │   │   │   showme-vprintf from doit: run-doit / false
2026/08/02 21:06:33 │   │   │   │   └ doit time: 3.026µs
2026/08/02 21:06:33 │   │   │   └ doit time: 13.035µs
2026/08/02 21:06:33 │   │   └ doit time: 27.08µs
2026/08/02 21:06:33 │   └ run
2026/08/02 21:06:33 │   ┌ doit: str: "main", recursive: false
2026/08/02 21:06:33 │   │   ┌ doit: no defer main false
2026/08/02 21:06:33 │   │   └ doit time: 1.435µs
2026/08/02 21:06:33 │   │   ┌ doit: fn main false
2026/08/02 21:06:33 │   │   │   printf from doit: main / false
2026/08/02 21:06:33 │   │   │   0-vprintf from doit: main / false
2026/08/02 21:06:33 │   │   └ doit time: 2.05µs
2026/08/02 21:06:33 │   └ doit time: 25.217µs
2026/08/02 21:06:33 │   exiting from main
2026/08/02 21:06:33 └ main
`)

	// Only all logs+times for doit when called thru run
	log.Reset()
	log.SetVerbose("run:doit:**,timed")
	// fmt.Printf("TEST: %s\n", log.Dump())
	main()
	diff(t, res, `
2026/08/02 21:12:38 printf from main
2026/08/02 21:12:38 0-vprintf from main
2026/08/02 21:12:38 ┌ run
2026/08/02 21:12:38 │   printf from run
2026/08/02 21:12:38 │   0-vprintf from run
2026/08/02 21:12:38 │   ┌ doit: str: "run", recursive: true
2026/08/02 21:12:38 │   │   ┌ doit: no defer run true
2026/08/02 21:12:38 │   │   └ doit time: 2.525µs
2026/08/02 21:12:38 │   │   ┌ doit: fn run true
2026/08/02 21:12:38 │   │   │   printf from doit: run / true
2026/08/02 21:12:38 │   │   │   0-vprintf from doit: run / true
2026/08/02 21:12:38 │   │   │   2-vprintf from doit: run / true
2026/08/02 21:12:38 │   │   │   showme-vprintf from doit: run / true
2026/08/02 21:12:38 │   │   └ doit time: 2.966µs
2026/08/02 21:12:38 │   │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 21:12:38 │   │   │   ┌ doit: no defer run-doit false
2026/08/02 21:12:38 │   │   │   └ doit time: 1.75µs
2026/08/02 21:12:38 │   │   │   ┌ doit: fn run-doit false
2026/08/02 21:12:38 │   │   │   │   printf from doit: run-doit / false
2026/08/02 21:12:38 │   │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 21:12:38 │   │   │   │   2-vprintf from doit: run-doit / false
2026/08/02 21:12:38 │   │   │   │   showme-vprintf from doit: run-doit / false
2026/08/02 21:12:38 │   │   │   └ doit time: 2.712µs
2026/08/02 21:12:38 │   │   └ doit time: 18.116µs
2026/08/02 21:12:38 │   └ doit time: 66.438µs
2026/08/02 21:12:38 └ run time: 71.946µs
2026/08/02 21:12:38 printf from doit: main / false
2026/08/02 21:12:38 0-vprintf from doit: main / false
2026/08/02 21:12:38 exiting from main
`)

	// Only all logs+times for doit when called thru run
	log.Reset()
	log.SetVerbose("run:doit:**,timed")
	// fmt.Printf("TEST: %s\n", log.Dump())
	main()
	diff(t, res, `
2026/08/02 21:13:08 printf from main
2026/08/02 21:13:08 0-vprintf from main
2026/08/02 21:13:08 ┌ run
2026/08/02 21:13:08 │   printf from run
2026/08/02 21:13:08 │   0-vprintf from run
2026/08/02 21:13:08 │   ┌ doit: str: "run", recursive: true
2026/08/02 21:13:08 │   │   ┌ doit: no defer run true
2026/08/02 21:13:08 │   │   └ doit time: 1.45µs
2026/08/02 21:13:08 │   │   ┌ doit: fn run true
2026/08/02 21:13:08 │   │   │   printf from doit: run / true
2026/08/02 21:13:08 │   │   │   0-vprintf from doit: run / true
2026/08/02 21:13:08 │   │   │   2-vprintf from doit: run / true
2026/08/02 21:13:08 │   │   │   showme-vprintf from doit: run / true
2026/08/02 21:13:08 │   │   └ doit time: 2.749µs
2026/08/02 21:13:08 │   │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 21:13:08 │   │   │   ┌ doit: no defer run-doit false
2026/08/02 21:13:08 │   │   │   └ doit time: 1.306µs
2026/08/02 21:13:08 │   │   │   ┌ doit: fn run-doit false
2026/08/02 21:13:08 │   │   │   │   printf from doit: run-doit / false
2026/08/02 21:13:08 │   │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 21:13:08 │   │   │   │   2-vprintf from doit: run-doit / false
2026/08/02 21:13:08 │   │   │   │   showme-vprintf from doit: run-doit / false
2026/08/02 21:13:08 │   │   │   └ doit time: 6.025µs
2026/08/02 21:13:08 │   │   └ doit time: 14.494µs
2026/08/02 21:13:08 │   └ doit time: 73.524µs
2026/08/02 21:13:08 └ run time: 77.092µs
2026/08/02 21:13:08 printf from doit: main / false
2026/08/02 21:13:08 0-vprintf from doit: main / false
2026/08/02 21:13:08 exiting from main
`)

	// Only 2's
	log.Reset()
	log.SetVerbose("2")
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS 2-vprintf from main
TS printf from run
TS 0-vprintf from run
TS 2-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS 2-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS 2-vprintf from doit: run-doit / false
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS 2-vprintf from doit: main / false
TS exiting from main
`)

	// Only showme's
	log.Reset()
	log.SetVerbose("showme")
	main()
	diff(t, res, `
TS printf from main
TS 0-vprintf from main
TS showme-vprintf from main
TS printf from run
TS 0-vprintf from run
TS showme-vprintf from run
TS printf from doit: run / true
TS 0-vprintf from doit: run / true
TS showme-vprintf from doit: run / true
TS printf from doit: run-doit / false
TS 0-vprintf from doit: run-doit / false
TS showme-vprintf from doit: run-doit / false
TS printf from doit: main / false
TS 0-vprintf from doit: main / false
TS showme-vprintf from doit: main / false
TS exiting from main
`)

	// Only showme's under doit
	log.Reset()
	log.SetVerbose("doit:showme")
	main()
	diff(t, res, `
2026/08/02 21:14:14 printf from main
2026/08/02 21:14:14 0-vprintf from main
2026/08/02 21:14:14 printf from run
2026/08/02 21:14:14 0-vprintf from run
2026/08/02 21:14:14 ┌ doit: str: "run", recursive: true
2026/08/02 21:14:14 │   ┌ doit: no defer run true
2026/08/02 21:14:14 │   └ doit
2026/08/02 21:14:14 │   ┌ doit: fn run true
2026/08/02 21:14:14 │   │   printf from doit: run / true
2026/08/02 21:14:14 │   │   0-vprintf from doit: run / true
2026/08/02 21:14:14 │   │   showme-vprintf from doit: run / true
2026/08/02 21:14:14 │   └ doit
2026/08/02 21:14:14 │   ┌ doit: str: "run-doit", recursive: false
2026/08/02 21:14:14 │   │   ┌ doit: no defer run-doit false
2026/08/02 21:14:14 │   │   └ doit
2026/08/02 21:14:14 │   │   ┌ doit: fn run-doit false
2026/08/02 21:14:14 │   │   │   printf from doit: run-doit / false
2026/08/02 21:14:14 │   │   │   0-vprintf from doit: run-doit / false
2026/08/02 21:14:14 │   │   │   showme-vprintf from doit: run-doit / false
2026/08/02 21:14:14 │   │   └ doit
2026/08/02 21:14:14 │   └ doit
2026/08/02 21:14:14 └ doit
2026/08/02 21:14:14 ┌ doit: str: "main", recursive: false
2026/08/02 21:14:14 │   ┌ doit: no defer main false
2026/08/02 21:14:14 │   └ doit
2026/08/02 21:14:14 │   ┌ doit: fn main false
2026/08/02 21:14:14 │   │   printf from doit: main / false
2026/08/02 21:14:14 │   │   0-vprintf from doit: main / false
2026/08/02 21:14:14 │   │   showme-vprintf from doit: main / false
2026/08/02 21:14:14 │   └ doit
2026/08/02 21:14:14 └ doit
2026/08/02 21:14:14 exiting from main
`)

}

func TestScoped(t *testing.T) {
	res := &strings.Builder{}
	golog.SetOutput(res)

	// doitScoped:2
	log.Reset()
	log.SetVerbose("doitScoped:2")
	mainScoped()
	diff(t, res, `
TS ┌ doitScoped
TS │   1-vprintf from doitScoped
TS │   1-vprintf from runScoped - from: doitScoped
TS │   2-vprintf from runScoped - from: doitScoped
TS └ doitScoped
`)

	// Just 1
	log.Reset()
	log.SetVerbose("1")
	mainScoped()
	diff(t, res, `
TS ┌ 1: 1-vprintf from mainScoped - enter
TS │   1-vprintf from mainScoped
TS │   1-vprintf from doitScoped
TS │   1-vprintf from runScoped - from: doitScoped
TS └ 1
TS 1-vprintf from runScoped - from: mainScoped
`)

	// 1+rr
	log.Reset()
	log.SetVerbose("1,rr")
	mainScoped()
	diff(t, res, `
TS ┌ 1: 1-vprintf from mainScoped - enter
TS │   1-vprintf from mainScoped
TS │   1-vprintf from doitScoped
TS │   ┌ rr: rr-v from runScoped - from: doitScoped
TS │   │   1-vprintf from runScoped - from: doitScoped
TS │   └ rr
TS └ 1
TS ┌ rr: rr-v from runScoped - from: mainScoped
TS │   1-vprintf from runScoped - from: mainScoped
TS └ rr
`)

	// rr:1   (nested, not & )
	log.Reset()
	log.SetVerbose("rr:1,mm")
	mainScoped()
	diff(t, res, `
TS mm-vprintf from mainScoped
TS ┌ rr: rr-v from runScoped - from: doitScoped
TS │   1-vprintf from runScoped - from: doitScoped
TS └ rr
TS ┌ rr: rr-v from runScoped - from: mainScoped
TS │   1-vprintf from runScoped - from: mainScoped
TS └ rr
`)

	// Same but no indent
	log.SetIndent(false)
	mainScoped()
	diff(t, res, `
TS mm-vprintf from mainScoped
TS > rr: rr-v from runScoped - from: doitScoped
TS 1-vprintf from runScoped - from: doitScoped
TS < rr
TS > rr: rr-v from runScoped - from: mainScoped
TS 1-vprintf from runScoped - from: mainScoped
TS < rr
`)

	// Indent back on but delete mm
	log.SetIndent(true)
	log.DelVerbose("mm")
	mainScoped()
	diff(t, res, `
TS ┌ rr: rr-v from runScoped - from: doitScoped
TS │   1-vprintf from runScoped - from: doitScoped
TS └ rr
TS ┌ rr: rr-v from runScoped - from: mainScoped
TS │   1-vprintf from runScoped - from: mainScoped
TS └ rr
`)

	// delete  0 at rr level
	// fmt.Printf("PRE VER: %s\n", log.String())
	log.DelVerbose("rr:0")
	// fmt.Printf("POST VER: %s\n", log.String())
	// log.DelVerbose("rr")
	// fmt.Printf("POST2 VER: %s\n", log.String())
	// fmt.Printf("POST2 DUMP: %s\n", log.Dump())
	mainScoped()
	diff(t, res, `
TS ┌ rr: rr-v from runScoped - from: doitScoped
TS └ rr
TS ┌ rr: rr-v from runScoped - from: mainScoped
TS └ rr
`)

	// trace w/args
	log.Reset()
	log.SetVerbose("doitScoped:runScoped:timed:**")
	// fmt.Printf("log:%s", log.Dump())
	mainScoped()
	diff(t, res, `
2026/08/02 21:16:24 ┌ doitScoped
2026/08/02 21:16:24 │   ┌ runScoped: from: doitScoped
2026/08/02 21:16:24 │   │   ┌ runScoped: doitScoped
2026/08/02 21:16:24 │   │   │   1-vprintf from runScoped - from: doitScoped
2026/08/02 21:16:24 │   │   │   2-vprintf from runScoped - from: doitScoped
2026/08/02 21:16:24 │   │   └ runScoped time: 10.561µs
2026/08/02 21:16:24 │   └ runScoped time: 61.831µs
2026/08/02 21:16:24 └ doitScoped
`)

}

func TestVerboseFunc(t *testing.T) {
	res := &strings.Builder{}
	golog.SetOutput(res)

	log.Reset()
	log.SetVerbose("mainIsFuncVerbose")
	mainIsFuncVerbose()
	diff(t, res, `2026/09/01 18:41:14 I'm only shown when mainVerboseFunc is on
2026/09/01 18:41:14 Me too!
`)
}
