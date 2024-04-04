FROM golang:1.22.2-alpine as go

FROM ubuntu:24.04

COPY --from=go /usr/local/go /usr/local/go
COPY --from=go /go /go

ENV GO_PATH="/usr/local/go"
ENV PATH="$PATH:$GO_PATH/bin"
ENV GO_ROOT="/go"
ENV PATH="$PATH:$GO_ROOT/bin"
ENV TZ="Asia/Ho_Chi_Minh"

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en

RUN apt update && apt upgrade -y
RUN apt install telnet -y && apt install curl -y && apt install vim -y && apt install nano -y
