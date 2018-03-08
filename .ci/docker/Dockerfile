FROM arm32v7/golang

RUN apt-get update


COPY selfhydro-release/selfhydro bin/

RUN /bin/selfhydro

ENTRYPOINT ["/selfhydro"]
