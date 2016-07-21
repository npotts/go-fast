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
	"github.com/alecthomas/kingpin"
	"github.com/alecthomas/units"
	"log"
	"os"

	"github.com/npotts/go-fast"
)

var (
	app     = kingpin.New("go-fast", "A CLI interface to www.fast.com")
	workers = app.Flag("workers", "Number of workers to start. Currently www.fast.com/Netflix only allows up to 3").Default("3").Short('w').Uint()
	bytes   = app.Flag("max", "Maximum worker download size").Default("50MB").Short('m').String()
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
	log.Printf("Starting with %d worker(s)\n", *workers)
	results := <-gf.Measure(int(*workers), int64(nb))
	log.Println(results)
}
