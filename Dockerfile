FROM golang:1.17 as build
COPY . /op-readiness
WORKDIR /op-readiness

RUN go build -o /op-readiness/main .
CMD ["/op-readiness/main"]
