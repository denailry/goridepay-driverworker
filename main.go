package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"goridepay-driverworker/common"
	"goridepay-driverworker/model"
	"goridepay-driverworker/worker"

	"github.com/gorilla/mux"
)

var port = os.Args[1]

func orderHandler(w http.ResponseWriter, r *http.Request) {
	orderRequest := model.NewOrderRequest(r.Body)
	orderInfo := model.OrderInfo{
		OrderID:     orderRequest.OrderID,
		Origin:      orderRequest.Origin,
		Destination: orderRequest.Destination,
	}
	for _, driverData := range orderRequest.DriverData {
		o := model.Order{
			Info:                &orderInfo,
			OriginDistance:      driverData.OriginDistance,
			DestinationDistance: driverData.DestinationDistance,
		}
		worker.AddOrder(driverData.DriverID, o)
	}
	or := model.OrderResponse{
		Error:   false,
		Message: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(or.ToJSON())
}

func main() {
	common.ServiceId, _ = strconv.Atoi(os.Args[2])

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/order", orderHandler).Methods("POST")
	http.Handle("/", r)
	// Add your routes as needed

	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Println("Listening on port " + port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Printf("\b\b")
	log.Println("shutting down")
	os.Exit(0)
}
