FROM arm32v7/ubuntu

ARG TAG
RUN echo $TAG

ADD https://github.com/bchalk101/selfhydro/releases/download/$TAG/selfhydro bin/selfhydro

RUN chmod 755 bin/selfhydro

CMD /bin/selfhydro&

