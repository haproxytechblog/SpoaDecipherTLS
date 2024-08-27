package main

import (
	"context"
	"fmt"
	ArgsParser "github.com/mmaFR/ArgumentsAsStruct"
	"github.com/mmaFR/signal_handler"
	"github.com/mmaFR/tls_handler"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/negasus/haproxy-spoe-go/request"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
)

var version string
var compileDate string
var commit string
var listener net.Listener
var logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
var workerCancel context.CancelFunc
var workerWg *sync.WaitGroup

func main() {
	const structure string = "core"
	var err error
	var spoa *agent.Agent
	var sigHandler *signal_handler.SignalHandler
	var args Arguments
	var nssKeylogFileDescriptor *os.File
	var worker *Worker
	var workerCtx context.Context

	ArgsParser.Parse(&args)

	if args.Version {
		fmt.Printf("Version: %s\nCompiled on: %s\nCommit: %s\n", version, compileDate, commit)
		return
	}

	if args.EnableTls {
		tls_handler.Logger = logger
	}
	if args.GenCa || args.GenSpoeCert || args.GenSpoaCert {
		tls_handler.Logger = logger
		tls_handler.GenerateCertificate(&args)
		return
	}

	args.LogOptions(logger)

	if nssKeylogFileDescriptor, err = os.OpenFile(args.NssKeylogFile, os.O_SYNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		logger.Fatalf("%s %s: error encountered while opening the NSS Keylog file: %s\n", structure, "main", err.Error())
	}
	defer nssKeylogFileDescriptor.Close()

	workerCtx, workerCancel = context.WithCancel(context.Background())
	worker = NewWorker(1024, nssKeylogFileDescriptor, workerCtx, logger)
	workerWg = worker.GetWg()
	go worker.Run()

	spoa = agent.New(SpoaHandler(worker.GetChannel()), NewLogger(logger))
	if args.EnableTls {
		if listener, err = tls_handler.NewListener(&args); err != nil {
			logger.Fatalf("%s %s: error encountered while creating the listener: %s\n", structure, "main", err.Error())
		}
	} else {
		if listener, err = net.Listen("tcp4", args.GetBindAddressAndPort()); err != nil {
			logger.Fatalf("%s %s: error encountered while creating the listener: %s\n", structure, "main", err.Error())
		}
	}
	defer listener.Close()
	sigHandler = signal_handler.NewSignalHandler(logger)
	_ = sigHandler.RegisterCallback(CallbackStop)
	if err = sigHandler.StartOn([]os.Signal{syscall.SIGINT, syscall.SIGTERM}); err != nil {
		logger.Printf("%s %s: error encountered while starting the signal handler: %s\n", structure, "main", err.Error())
	}

	if err = spoa.Serve(listener); err != nil {
		logger.Printf("%s %s: error encountered while starting the spoa server: %s\n", structure, "main", err.Error())
	}

	sigHandler.Wait()
}

func SpoaHandler(reqChan chan *SpopRequest) func(req *request.Request) {
	return func(req *request.Request) {
		var wrappedRequest *SpopRequest
		var reqCtx context.Context
		logger.Println("Core SpoaHandler: new request received")
		wrappedRequest, reqCtx = NewRequestWithCancel(req)
		reqChan <- wrappedRequest
		<-reqCtx.Done()
	}
}

func CallbackStop(_ os.Signal) {
	logger.Println("SignalHandler CallbackStop: trying to close the listener")
	_ = listener.Close()
	logger.Println("SignalHandler CallbackStop: listener closed")
	logger.Println("SignalHandler CallbackStop: trying to stop the worker")
	workerCancel()
	workerWg.Wait()
	logger.Println("SignalHandler CallbackStop: worker stopped")
}
