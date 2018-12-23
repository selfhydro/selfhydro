FROM alpine:latest

RUN apk update && pk --no-cache add mosquitto

EXPOSE 1883

RUN mkdir -p /mqtt/

COPY mosquitto/mosquitto.conf /mqtt/config/

RUN mkdir -p /var/lib/mosquitto/

RUN adduser --system --disabled-password --disabled-login mosquitto

RUN /sbin/ldconfig

ENTRYPOINT ["/usr/sbin/mosquitto", "-c", "/etc/mosquitto/mosquitto.conf"]

