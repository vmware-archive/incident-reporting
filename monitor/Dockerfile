FROM golang:1.11.4

WORKDIR /go/src/app
COPY *.go .

RUN go get -d -v ./... && \
    CGO_ENABLED=0 GOOS=linux go build -v


FROM scratch
WORKDIR /
COPY --from=0 /go/src/app/app .
CMD ["./app"]