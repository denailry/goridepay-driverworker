FROM golang:latest

RUN go get -u github.com/gorilla/mux

ADD main.go /go/src/goridepay-driverworker/main.go
ADD common /go/src/goridepay-driverworker/common
ADD worker /go/src/goridepay-driverworker/worker

CMD go run /go/src/goridepay-driverworker/main.go $PORT