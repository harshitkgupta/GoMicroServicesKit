FROM golang:1.13 AS builder
WORKDIR /app
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/server .
ENTRYPOINT [ "/root/server" ]