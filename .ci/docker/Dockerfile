FROM arm32v6/alpine

ARG TAG

USER selfhydro

RUN apt update && apt install dumb-init

RUN mkdir -p /bin/selfhydro \
    && curl -SL https://github.com/selfhydro/selfhydro/releases/download/$TAG/selfhydro

RUN chmod 755 bin/selfhydro/selfhydro

ENTRYPOINT ["dumbinit", "/bin/selfhydro"]
