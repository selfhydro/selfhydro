FROM debian

RUN apt-get update -y
RUN apt-get install -y wget

ARG TAG

RUN wget https://github.com/bchalk101/selfhydro/releases/download/$TAG/selfhydro

COPY selfhydro bin/

CMD /bin/selfhydro

ENTRYPOINT ["/selfhydro"]
