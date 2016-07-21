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
	Error    error         `json:"error"`             //Error, if any, that occurred while reading in a sample set
	Duration time.Duration `json:"duation"`           //how long the operation took to read in Bytes
	Bytes    int           `json:"bytes"`             //the number of bytes read in
	Bps      float64       `json:"bps"`               //bits per second over the immediate sample period
	Samples  []Stats       `json:"samples,omitempty"` //Workers should populate this with all their stats
}

func (s Stats) String() string {
	return fmt.Sprintf("Error=%v Duration=%v # Bytes=%d Bps=%4.3f", s.Error, s.Duration, s.Bytes, s.Bps)
}

//nStats is a slice of Stats
type nStats []Stats

func (n nStats) Stats() (rtn Stats) {
	rtn.Samples = n //store all individual samples here
	for _, i := range n {
		rtn.Bytes += i.Bytes
		rtn.Duration += i.Duration
	}
	rtn.Bps = bps(rtn.Duration, rtn.Bytes)
	return
}

//Results is the final results of the test complete with raw data measured form each worker
type Results struct {
	Bytes      []int           `json:"bytes,omitempty"`      //A slice of the total bytes read by each worker
	Duration   []time.Duration `json:"duration,omitempty"`   // slice of total duration the operation took
	BitsPerSec []float64       `json:"bitspersec,omitempty"` //sice of calculated bits/sec from each worker
	Samples    []nStats        `json:"samples,omitempty"`    //Slice of all the samples each worker measured.
	Workers    int             `json:"workers"`              //number of workers used
	Bps        float64         `json:"bps"`                  //bits per second
	Kbps       float64         `json:"kbps"`                 //kbits per second
	Mbps       float64         `json:"mbps"`                 //MBits per second
}

func (r Results) String() string {
	return fmt.Sprintf(`%d worker(s) downloaded at an average of
  %.2f Bps
  %.2f Kbps
  %.2f Mbps`, r.Workers, r.Bps, r.Kbps, r.Mbps)
}
