FROM golang:1.22 AS Builder

WORKDIR /build

COPY ./* ./*