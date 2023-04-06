package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"www.hello.com/services/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodBye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChain := make(chan os.Signal)
	signal.Notify(sigChain, os.Interrupt)
	signal.Notify(sigChain, os.Kill)

	sig := <-sigChain
	l.Println("Recieved terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

func server(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL Path" + r.URL.Path)
	fmt.Fprintf(w, "Hello World %s", r.URL.Path[1:])
}
