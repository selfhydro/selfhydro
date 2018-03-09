FROM debian

RUN apt-get -y update
RUN apt-get -y install crossbuild-essential-armhf

ARG TAG

ADD https://github.com/bchalk101/selfhydro/releases/download/$TAG/selfhydro bin/selfhydro

RUN chmod 755 bin/selfhydro

CMD /bin/selfhydro&

