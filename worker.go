package gofast

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
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
	ID    int
}

func (w *worker) Stat() Stats {
	return w.stats.Stats()
}

func (w *worker) Start(url string, cfg Settings, wg *sync.WaitGroup) {
	defer wg.Done()
	tlast := time.Now()
	tstart := time.Now()
	total := int64(0)
	n := 0
	var err error
	resp := &http.Response{}

	if resp, err = http.Get(url); err != nil {
		err = errors.Wrapf(err, "Worker % 2d unable to initalize GET to %s", w.ID, url)
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	for {
		if n, err = w.read(resp.Body, cfg.Network); err != nil {
			err = errors.Wrapf(err, "Worker % 2d encounted timeout", w.ID)
			log.Println(err)
			return
		}
		nstat := Stats{Duration: time.Since(tlast), Bytes: n, Error: err}
		tlast = time.Now()
		nstat.Bps = bps(nstat.Duration, nstat.Bytes)
		w.stats = append(w.stats, nstat)
		if total += int64(n); (cfg.MaxBytes > 0 && total > cfg.MaxBytes) ||
			(cfg.Timeout > 0 && time.Since(tstart) > cfg.Timeout) {
			return
		}
	}
}

func (w *worker) read(reader io.Reader, timeout time.Duration) (n int, err error) {
	p := make([]byte, 1024*32)
	if timeout > 0 {
		ch := make(chan bool)
		go func() {
			n, err = reader.Read(p)
			ch <- true
		}()
		select {
		case <-ch:
			return
		case <-time.After(timeout):
			err := fmt.Errorf("Read() on socket timed out after %v", timeout)
			log.Println(err)
			return 0, err
		}
	} else {
		return reader.Read(p)
	}

}
