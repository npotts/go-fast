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

package gofast

import (
	"fmt"
	"time"
)

func bps(dur time.Duration, bytes int) float64 {
	return 8 * float64(bytes) / dur.Seconds()
}

//Stats is a simple store of statistic data
type Stats struct {
	Error    error
	Duration time.Duration
	Bytes    int
	Bps      float64 //bits per second
}

func (s Stats) String() string {
	return fmt.Sprintf("Error=%v Duration=%v # Bytes=%d Bps=%4.3f", s.Error, s.Duration, s.Bytes, s.Bps)
}

type nStats []Stats

func (n nStats) Stats() (rtn Stats) {
	for _, i := range n {
		rtn.Bytes += i.Bytes
		rtn.Duration += i.Duration
	}
	rtn.Bps = bps(rtn.Duration, rtn.Bytes)
	return
}

//Results is the final results of the test
type Results struct {
	Bytes      []int
	Duration   []time.Duration
	BitsPerSec []float64
	Workers    int
	Bps        float64
	Kbps       float64
	Mbps       float64
}

func (r Results) String() string {
	return fmt.Sprintf(`%d worker(s) downloaded at an average of
  %.2f Bps
  %.2f Kbps
  %.2f Mbps`, r.Workers, r.Bps, r.Kbps, r.Mbps)
}
