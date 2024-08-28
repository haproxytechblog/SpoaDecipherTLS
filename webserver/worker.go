package main

import (
	"DecipherTLS/cache"
	"context"
	"log"
	"os"
	"sync"
)

type Worker struct {
	fd        *os.File
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc
	reqChan   chan []byte
	cache     *cache.Cache
	logger    *log.Logger
}

func (w *Worker) GetRequestChan() chan<- []byte {
	return w.reqChan
}
func (w *Worker) Wait() {
	w.wg.Wait()
}
func (w *Worker) Stop() {
	w.ctxCancel()
}
func (w *Worker) Run() {
	var err error
	var loop bool = true
	var request []byte
	var tlsData *TlsData
	w.wg.Add(1)
	w.logger.Println("Worker starts")
	for loop {
		select {
		case request = <-w.reqChan:
			if tlsData, err = NewTlsData(request); err != nil {
				w.logger.Println(err)
				continue
			}
			w.logger.Printf("Request received for the client random %s, tls version used %s\n", tlsData.ClientRandom, tlsData.ProtocolVersion)
			if !w.cache.Exists(tlsData.Sprint()) {
				if _, err = w.fd.WriteString(tlsData.Sprint()); err != nil {
					w.logger.Printf("Worker Run: failed to write tls data: %s\n", err)
				}
				if err = w.fd.Sync(); err != nil {
					w.logger.Printf("Worker Run: failed to sync: %s\n", err)
				}
			} else {
				w.logger.Println("Worker Run: tls data already exists in cache, so ignored")
			}
		case <-w.ctx.Done():
			if err = w.fd.Sync(); err != nil {
				w.logger.Println(err)
			}
			if err = w.fd.Close(); err != nil {
				w.logger.Println(err)
			}
			loop = false
		}
	}
	w.logger.Println("Worker stops")
	w.wg.Done()
}

func NewWorker(filename string, cacheSize uint) (*Worker, error) {
	var err error
	var w *Worker = new(Worker)
	var fd *os.File
	if fd, err = os.OpenFile(filename, os.O_SYNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		return nil, err
	}
	w.fd = fd
	w.ctx, w.ctxCancel = context.WithCancel(context.Background())
	w.reqChan = make(chan []byte)
	w.cache = cache.New(cacheSize)
	w.logger = logger
	return w, nil
}
