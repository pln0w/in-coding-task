FROM golang:alpine AS builder

ENV CGO_ENABLED=1

COPY . /go/src/server
WORKDIR /go/src/server

RUN apk add --no-cache git

RUN go get -v ./...
RUN go install -v ./...

RUN go build -o server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/server /server

ENTRYPOINT ./server

LABEL Name=server Version=0.0.1

EXPOSE 8080