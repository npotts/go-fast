package gofast

import (
	"fmt"
	// "io"
	"net/http"
	"sync"
	"time"
)

var _ = fmt.Sprintf

/*The Worker interface is used by the actual lower level workers
responsible for downloading the files*/
type Worker interface {
	Start(url string, maxsize int, wg *sync.WaitGroup)
	Stat() Stats
}

/*worker does the job of downloading the data located at the passed file and */
type worker struct {
	url   string
	stats nStats
}

func (w *worker) Stat() Stats {
	return w.stats.Stats()
}

func (w *worker) Start(url string, maxsize int, wg *sync.WaitGroup) {
	fmt.Println("Worker fetching from ", url)
	tlast := time.Now()
	total := 0
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		for {
			p := make([]byte, 1024*16)
			n, e := resp.Body.Read(p)
			nstat := Stats{Duration: time.Since(tlast), Bytes: n, Error: e}
			tlast = time.Now()
			nstat.Bps = bps(nstat.Duration, nstat.Bytes)
			w.stats = append(w.stats, nstat)
			if total += n; e != nil || total > maxsize {
				break
			}
		}
	}
	wg.Done()
}
