FROM golang:1.18 as build
WORKDIR /go/src/github.com/k8sbykeshed/op-readiness
COPY . .
RUN curl -L "https://dl.k8s.io/v1.23.5/kubernetes-test-linux-amd64.tar.gz" -o /tmp/test.tar.gz
RUN tar xvzf /tmp/test.tar.gz --strip-components=3 kubernetes/test/bin/e2e.test
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o op-readiness .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=0 /go/src/github.com/k8sbykeshed/op-readiness/e2e.test /app/
COPY --from=0 /go/src/github.com/k8sbykeshed/op-readiness/op-readiness /app/
COPY --from=0 /go/src/github.com/k8sbykeshed/op-readiness/tests.yaml /app/
CMD ["./op-readiness"]
