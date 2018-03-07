FROM arm32v7/golang

RUN apt-get update

COPY selfhydro bin/

RUN ./bin/selfhydro

ENTRYPOINT ["/selfhydro"]
