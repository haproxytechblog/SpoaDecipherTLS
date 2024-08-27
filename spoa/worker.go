package main

import (
	"DecipherTLS/cache"
	"context"
	"fmt"
	"github.com/negasus/haproxy-spoe-go/action"
	message2 "github.com/negasus/haproxy-spoe-go/message"
	"log"
	"os"
	"sync"
)

type Worker struct {
	cache   *cache.Cache
	reqChan chan *SpopRequest
	fd      *os.File
	ctx     context.Context
	wg      *sync.WaitGroup
	logger  *log.Logger
}

func (w *Worker) GetChannel() chan *SpopRequest {
	return w.reqChan
}

func (w *Worker) GetWg() *sync.WaitGroup {
	return w.wg
}

func (w *Worker) Run() {
	w.wg.Add(1)
	defer w.wg.Done()
	var err error
	var req *SpopRequest
	var tlsData *TlsData
	var message *message2.Message
	var loop bool = true

	w.logger.Println("Worker Run: started")

	for loop {
		select {
		case req = <-w.reqChan:
			req.Actions.SetVar(action.ScopeTransaction, "OK", true)
			for _, msgName := range []string{"fc_ssl_params", "bc_ssl_params"} {
				message, err = req.Messages.GetByName(msgName)
				switch {
				case err == nil:
					if tlsData, err = NewTlsData(message); err != nil {
						w.logger.Printf("Worker Run: failed to parse tls data: %s\n", err)
						req.Done()
						continue
					}
					if !w.cache.Exists(tlsData.Sprint()) {
						var i int
						if i, err = w.fd.WriteString(tlsData.Sprint()); err != nil {
							w.logger.Printf("Worker Run: failed to write tls data: %s\n", err)
						}
						fmt.Println(tlsData.Sprint())
						fmt.Printf("%d bytes wrote in the nss file\n", i)
						if err = w.fd.Sync(); err != nil {
							w.logger.Printf("Worker Run: failed to sync: %s\n", err)
						}
					} else {
						w.logger.Println("Worker Run: tls data already exists in cache, so ignored")
					}
				case err.Error() == "message not found":
				default:
					w.logger.Printf("Worker Run: get message error: %s\n", err)
				}
			}
			req.Done()
		case <-w.ctx.Done():
			loop = false
			w.logger.Println("Worker Run: exiting")
			continue
		}
	}
}

func NewWorker(cacheSize uint, fd *os.File, ctx context.Context, logger *log.Logger) *Worker {
	var w *Worker = new(Worker)
	w.cache = cache.New(cacheSize)
	w.reqChan = make(chan *SpopRequest)
	w.fd = fd
	w.ctx = ctx
	w.wg = new(sync.WaitGroup)
	w.logger = logger
	return w
}
