FROM golang:1.20-alpine AS builder

WORKDIR /build
COPY . .

RUN GOPROXY=https://goproxy.cn go build -o wake ./cmd/wake

FROM alpine
COPY --from=builder /build/wake /bin/wake

CMD ["wake"]
