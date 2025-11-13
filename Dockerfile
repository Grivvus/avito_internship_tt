FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server /server

ENTRYPOINT ["/server"]
