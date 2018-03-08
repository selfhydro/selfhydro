FROM debian

RUN apt-get update -y
RUN apt-get install -y wget

ARG TAG

RUN wget https://github.com/bchalk101/selfhydro/releases/download/$TAG/selfhydro

RUN mv selfhydro bin/selfhydro

CMD /bin/selfhydro&

ENTRYPOINT ["/selfhydro"]
