FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY . .
RUN go build -o main ./watermark/ \
    && apk --update add ca-certificates

FROM scratch
EXPOSE 80
COPY --from=builder /build/main /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["./main"]
