/*Package gofast is a go module that access www.fast.com in order to derive upload/download speeds*/
package main

/*
Copyright (c) 2016 Nicholas Potts

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without
limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/alecthomas/units"
	"log"
	"os"

	"github.com/npotts/go-fast"
)

var (
	app      = kingpin.New("go-fast", "A CLI interface to www.fast.com")
	workers  = app.Flag("workers", "Number of workers to start. Currently www.fast.com/Netflix only allows up to 3").Default("3").Short('w').Uint()
	bytes    = app.Flag("max", "Maximum worker download size. Default of 0 means to download the entirity of the files").Default("0").Short('m').String()
	network  = app.Flag("network", "Network timeout Default of should be sufficient").Default("5s").Short('n').Duration()
	timeout  = app.Flag("timeout", "Maximum time to allow workers to run. Default of 0s indicates to never timeout").Default("0s").Short('t').Duration()
	emitjson = app.Flag("json", "emit raw json data with all samples.  Implies --quiet").Default("false").Short('j').Bool()
	quiet    = app.Flag("quiet", "Only emit measured Bits per second value (and fatal errors)").Default("false").Short('q').Bool()
)

func parse() {

}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	nb, err := units.ParseBase2Bytes(*bytes)
	if err != nil {
		log.Printf("Invalid max size selection %q: %v", *bytes, err)
		os.Exit(-1)
	}
	if *workers == 0 {
		log.Println("Workers must be greater than zero")
		os.Exit(-1)
	}

	gf := gofast.New()
	cfg := gofast.Settings{MaxBytes: int64(nb), Timeout: *timeout, Workers: int(*workers), Network: *network}
	if !*quiet && !*emitjson {
		log.Printf("Starting with %d worker(s)\n", *workers)
	}
	results := <-gf.Measure(cfg)
	if results.Workers == 0 {
		os.Exit(-1)
	}

	if *emitjson {
		d, err := json.Marshal(results)
		fmt.Println(string(d))
		if err == nil {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if *quiet {
		fmt.Println(results.Bps)
		os.Exit(0)
	}
	log.Println(results)
}
