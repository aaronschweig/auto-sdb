FROM golang:1.17-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o main .


FROM alpine:latest

RUN apk add ghostscript

COPY --from=builder /build/main /

CMD ["/main"]
