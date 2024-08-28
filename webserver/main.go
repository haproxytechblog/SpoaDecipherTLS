package main

import (
	"fmt"
	ArgsParser "github.com/mmaFR/ArgumentsAsStruct"
	"github.com/mmaFR/signal_handler"
	"github.com/mmaFR/tls_handler"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

const urlPath string = "/newdata"

var version string
var compileDate string
var commit string
var logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
var worker *Worker
var listener net.Listener
var workerChan chan<- []byte

func main() {
	const structure string = "core"
	var err error
	var sigHandler *signal_handler.SignalHandler
	var args Arguments
	var mux *http.ServeMux = http.NewServeMux()

	ArgsParser.Parse(&args)

	if args.Version {
		fmt.Printf("Version: %s\nCompiled on: %s\nCommit: %s\n", version, compileDate, commit)
		return
	}

	if args.EnableTls {
		tls_handler.Logger = logger
	}
	if args.GenCa || args.GenHaproxyCert || args.GenServerCert {
		tls_handler.Logger = logger
		tls_handler.GenerateCertificate(&args)
		return
	}

	args.LogOptions(logger)

	if worker, err = NewWorker(args.NssKeylogFile, 1024); err != nil {
		logger.Fatalf("%s %s: error encountered while creating the worker: %s\n", structure, "main", err.Error())
	}
	workerChan = worker.GetRequestChan()

	go worker.Run()

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

	mux.HandleFunc(urlPath, RequestHandler)
	if err = http.Serve(listener, mux); err != nil {
		logger.Printf("%s %s: error encountered with server: %s\n", structure, "main", err.Error())
	}

	sigHandler.Wait()
	logger.Println("Bye bye !!")
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var payload []byte

	if payload, err = io.ReadAll(r.Body); err != nil {
		logger.Printf("core RequestHandler: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte{}); err != nil {
			logger.Printf("core RequestHandler: %s\n", err.Error())
		}
	} else {
		workerChan <- payload
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write([]byte{}); err != nil {
			logger.Printf("core RequestHandler: %s\n", err.Error())
		}
	}
}

func CallbackStop(_ os.Signal) {
	logger.Println("SignalHandler CallbackStop: trying to close the server")
	_ = listener.Close()
	logger.Println("SignalHandler CallbackStop: server closed")
	logger.Println("SignalHandler CallbackStop: trying to stop the worker")
	worker.Stop()
	worker.Wait()
	logger.Println("SignalHandler CallbackStop: worker stopped")
}
