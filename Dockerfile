FROM golang:1.11.2-alpine

RUN apt update && apt get git
RUN apk update && apk add git -y
