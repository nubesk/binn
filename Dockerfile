FROM golang:1.17.2-alpine3.14 AS builder
ENV GO111MODULE=on

RUN apk --update add alpine-sdk

WORKDIR /go/src/github.com/binn
COPY ./ ./
RUN go mod download && \
    go get github.com/binn/binn
RUN go build -o /tmp/server ./main.go


FROM alpine:latest
COPY --from=builder /tmp/server .

ENV BINN_SEED=42 \
    BINN_DELIVERY_CYCLE_SEC=20 \
    BINN_ENABLE_VALIDATION=true \
    BINN_GENERATE_CYCLE_SEC=60 \
    BINN_SEND_EMPTY_SEC=29 \
    BINN_ENGINE_ENABLE_DEBUG=false \
    BINN_SERVER_ENABLE_DEBUG=false

CMD ["./server"]
