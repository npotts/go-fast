package gofast

import (
	"net/http"
	"sync"
	"time"
)

/*The Worker interface is used by the actual lower level workers
responsible for downloading the files*/
type Worker interface {
	Start(url string, cfg Settings, wg *sync.WaitGroup)
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

func (w *worker) Start(url string, cfg Settings, wg *sync.WaitGroup) {
	tlast := time.Now()
	tstart := time.Now()
	total := int64(0)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		for {
			p := make([]byte, 1024*16)
			n, e := resp.Body.Read(p)
			nstat := Stats{Duration: time.Since(tlast), Bytes: n, Error: e}
			tlast = time.Now()
			nstat.Bps = bps(nstat.Duration, nstat.Bytes)
			w.stats = append(w.stats, nstat)
			if total += int64(n); e != nil ||
				(cfg.MaxBytes > 0 && total > cfg.MaxBytes) ||
				(cfg.Timeout > 0 && time.Since(tstart) > cfg.Timeout) {
				break
			}
		}
	}
	wg.Done()
}
