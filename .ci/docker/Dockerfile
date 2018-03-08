FROM arm32v7/golang:1.8

RUN apt-get update

RUN apt-get install -y wget

ARG TAG

RUN wget https://github.com/bchalk101/selfhydro/releases/download/v$TAG/selfhydro

COPY selfhydro bin/

RUN /bin/selfhydro

ENTRYPOINT ["/selfhydro"]
