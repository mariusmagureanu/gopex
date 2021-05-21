FROM golang:alpine AS builder
LABEL maintainer="marius@archlinux.live"
LABEL version="0.1"
ENV GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS=linux
COPY . /go/
WORKDIR /go/
RUN apk add --update make git
RUN make all

FROM alpine
RUN adduser -D -g '' kc
USER kc
COPY --from=builder /go/src/bin/monitor monitor
