package main

import (
	storage "RuCache/Storage"
	"RuCache/handler"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

func main() {
	var isShuttingDown atomic.Bool
	var err error
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	store := storage.NewStore()
	mux := http.NewServeMux()
	h := handler.NewHandler(store, isShuttingDown.Load)
	h.RegisterHandlers(mux)

	//
	//
	// test code
	//
	//

	ttl := 15 * time.Second
	// stats, err := store.Set("test1", "value1", nil)
	// stats, err = store.Set("test2", "value1", nil)
	// stats, err = store.Set("test3", "value1", nil)
	// stats, err = store.Set("test4", "value1", nil)
	// stats, err = store.Set("test5", "value1", nil)
	// stats, err = store.Set("test6", "value1", nil)
	stats, err := store.Set("test999", "value1", &ttl)
	if !stats {
		fmt.Println(err)
	}
	// fmt.Println(store.Get("test1"))
	// stats, err = store.Del("test1")
	// if !stats {
	// 	fmt.Println(err)
	// }
	fmt.Println(store.Get("test5"))

	//
	//
	// test code
	//
	//

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr: ":8080",
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
		Handler: mux,
	}

	go func() {
		log.Print("Server start on :" + server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Print("Received shutdown signal, waiting for ongoing requests to finish...")

	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err = server.Shutdown(shutDownCtx)
	err = store.SaveDataToFile()
	if err != nil {
		log.Printf("Save data to file failed: %v", err)
	}
	stopOngoingGracefully()
	if err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
		time.Sleep(_shutdownHardPeriod)
	}

	fmt.Println("Shutting down gracefully...")
	os.Exit(0)
}
