FROM golang:1.17-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .


FROM alpine:latest

RUN apk add ghostscript

COPY --from=builder /dist/main /

CMD ["/main"]
