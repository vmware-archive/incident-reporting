FROM golang:1.11.4

WORKDIR /go/src/app
RUN apt-get update && \
    apt-get install -y xz-utils && \
    wget https://github.com/ethereum/solidity/releases/download/v0.4.24/solc-static-linux -O /bin/solc && \
    chmod 755 /bin/solc

COPY ui-server/*.go /go/src/app/
COPY ui-server/public/ /go/src/app/public/
COPY truffle/contracts /go/src/app/contracts

RUN go get -d -v ./... && \
  cd /go/src/github.com/ethereum/go-ethereum && \
  go install ./... && \
  cd /go/src/app && \
  mkdir generate && \
  go generate -v && \
  CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v

EXPOSE 80

FROM busybox:glibc
WORKDIR /
COPY --from=0 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /go/src/app .
CMD ["./app"]