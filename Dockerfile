FROM golang:1.19 as build
WORKDIR /go/src/sigs.k8s.io/windows-operational-readiness
COPY . .
ARG KUBERNETES_VERSION
RUN curl -L "https://dl.k8s.io/${KUBERNETES_VERSION}/kubernetes-test-linux-amd64.tar.gz" -o /tmp/test.tar.gz
RUN tar xvzf /tmp/test.tar.gz --strip-components=3 kubernetes/test/bin/e2e.test && rm -f /tmp/test.tar.gz
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o op-readiness .

FROM debian:bookworm-slim
WORKDIR /app
ENV ARTIFACTS /tmp/sonobuoy/results
COPY --from=0 /go/src/sigs.k8s.io/windows-operational-readiness/e2e.test /app/
COPY --from=0 /go/src/sigs.k8s.io/windows-operational-readiness/op-readiness /app/
COPY --from=0 /go/src/sigs.k8s.io/windows-operational-readiness/tests.yaml /app/
CMD ["./op-readiness"]
