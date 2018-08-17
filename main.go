package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	"goridepay-driverworker/common"
	"goridepay-driverworker/invalidator"
	"goridepay-driverworker/model/accept"
	"goridepay-driverworker/model/common"
	"goridepay-driverworker/model/invalidate"
	"goridepay-driverworker/model/order"
	"goridepay-driverworker/model/reject"
	"goridepay-driverworker/worker"

	"github.com/gorilla/mux"
)

var port = os.Args[1]

func orderHandler(w http.ResponseWriter, r *http.Request) {
	request := order.NewRequest(r.Body)
	info := order.Info{
		OrderID:     request.OrderID,
		Origin:      request.Origin,
		Destination: request.Destination,
	}
	for _, driverData := range request.DriverData {
		o := order.Order{
			Info:                &info,
			OriginDistance:      driverData.OriginDistance,
			DestinationDistance: driverData.DestinationDistance,
		}
		worker.AddOrder(driverData.DriverID, o)
	}
	response := response.Response{
		Error:   false,
		Message: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response.ToJSON())
}

func acceptHandler(w http.ResponseWriter, r *http.Request) {
	request := accept.NewRequest(r.Body)
	var re response.Response
	if worker.AcceptOrder(request.DriverID, request.OrderID) {
		re.Error = false
		re.Message = "ok"
	} else {
		re.Error = true
		re.Message = "Order has been taken or cancelled."
	}
	w.WriteHeader(http.StatusOK)
	w.Write(re.ToJSON())
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {

}

func invalidateHandler(w http.ResponseWriter, r *http.Request) {
	request := invalidate.NewRequest(r.Body)
	invalidOrder := invalidate.NewInvalidOrder(request.OrderID)
	go invalidator.Invalidate(invalidOrder)
	re := response.Response{
		Error:   false,
		Message: "ok",
	}
	w.WriteHeader(http.StatusOK)
	w.Write(re.ToJSON())
}

func rejectHandler(w http.ResponseWriter, r *http.Request) {
	request := reject.NewRequest(r.Body)
	var re response.Response
	if worker.RejectOrder(request.DriverID, request.OrderID) {
		re.Error = false
		re.Message = "ok"
	} else {
		re.Error = true
		re.Message = "Order ID is not valid to be rejected."
	}
	w.WriteHeader(http.StatusOK)
	w.Write(re.ToJSON())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	runtime.GOMAXPROCS(128)
	common.ServiceId, _ = strconv.Atoi(os.Args[2])

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/order", orderHandler).Methods("POST")
	r.HandleFunc("/accept", acceptHandler).Methods("POST")
	r.HandleFunc("/cancel", cancelHandler).Methods("POST")
	r.HandleFunc("/reject", rejectHandler).Methods("POST")
	r.HandleFunc("/invalidate", rejectHandler).Methods("POST")
	r.HandleFunc("/", homeHandler).Methods("GET")
	http.Handle("/", r)

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
