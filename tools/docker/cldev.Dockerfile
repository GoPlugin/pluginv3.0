FROM golang:1.21-bullseye

ARG SRCROOT=/usr/local/src/plugin
WORKDIR ${SRCROOT}

ADD go.* ./
RUN go mod download
RUN mkdir -p tools/bin
