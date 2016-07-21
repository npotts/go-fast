/*Package gofast is a go module that accesses www.fast.com in order to derive upload/download speeds*/
package gofast

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
	"sync"
	"time"
)

/*Settings is the configurable settings to pass to a Measurer*/
type Settings struct {
	Workers  int           //number of workers to use
	MaxBytes int64         //maximum number of bytes to ingest
	Timeout  time.Duration //maximal time to run for
	EmitJSON bool          //emit JSON output
}

//Measurer is an interface used to measure values.  The returned channel will be written to exactly once
type Measurer interface {
	Measure(Settings) <-chan Results
}

//New returns an object that is a measurer
func New() Measurer {
	return new(gofast)
}

//basic structure that implements the Measurer interface
type gofast struct {
	token    string
	routines int
	stats    chan Results
}

//Measure implemented the measurement interface as well as performs the measurements
func (gf *gofast) Measure(cfg Settings) <-chan Results {
	urls, err := gf.getURLs(cfg.Workers)
	if err != nil {
		panic(err)
	}
	gf.stats = make(chan Results)
	if len(urls) == 0 {
		go func() { gf.stats <- Results{} }()
	} else {
		go gf.run(urls, cfg)
	}
	return gf.stats
}

func (gf *gofast) run(urls []string, cfg Settings) {
	//TODO: Fan-out to run tests, fan in with results
	var wg sync.WaitGroup
	workers := []Worker{}
	for _, url := range urls {
		wg.Add(1)
		worker := new(worker)
		go worker.Start(url, cfg, &wg)
		workers = append(workers, worker)
	}
	wg.Wait()
	stats := Results{Bytes: []int{}, Duration: []time.Duration{}, BitsPerSec: []float64{}, Workers: len(workers), Samples: []nStats{}}
	for _, worker := range workers {
		wstat := worker.Stat()
		stats.Bytes = append(stats.Bytes, wstat.Bytes)
		stats.Duration = append(stats.Duration, wstat.Duration)
		stats.BitsPerSec = append(stats.BitsPerSec, wstat.Bps)
		stats.Samples = append(stats.Samples, wstat.Samples)
		stats.Bps += wstat.Bps
	}
	stats.Kbps = stats.Bps / 1024.0
	stats.Mbps = stats.Kbps / 1024.0
	gf.stats <- stats
}
