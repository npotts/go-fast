/*Package gofast is a go module that access www.fast.com in order to derive upload/download speeds*/
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

import "time"

//Stats is a simple store of statistic data
type Stats struct {
	Duration time.Duration
	Bytes    int
	Kbps     float64
}

//Measurer is an interface used to measure values.  The returned channel will be written to exactly once
type Measurer interface {
	Measure(int) <-chan []Stats
}

//basic structure that implements the Measurer interface
type gofast struct {
	token    string
	routines int
	stats    chan []Stats
}

//Measure implemented the measurement interface as well as performs the measurements
func (gf *gofast) Measure(count int) (stats chan []Stats) {
	urls, err := gf.getURLs(count)
	if err != nil {
		panic(err)
	}
	gf.stats = make(chan []Stats, len(urls))
	if len(urls) == 0 {
		go func() { gf.stats <- []Stats{} }()
	} else {
		go gf.fanout()
	}
	return
}

func (gf *gofast) fanout() {
	//TODO: Fan-out to run tests, fan in with results

}
