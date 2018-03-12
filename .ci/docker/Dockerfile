FROM armhf/debian

ARG TAG

ADD https://github.com/bchalk101/selfhydro/releases/download/$TAG/selfhydro bin/selfhydro

RUN chmod 755 bin/selfhydro

CMD /bin/selfhydro&

