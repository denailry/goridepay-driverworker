FROM golang:latest

RUN go get -u github.com/gorilla/mux

ADD main.go /go/src/goridepay-driverworker/main.go

CMD go run /go/src/goridepay-driverworker/main.go $PORT