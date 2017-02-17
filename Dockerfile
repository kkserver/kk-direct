FROM alpine:latest

COPY ./etc/timezone /etc/timezone

COPY ./etc/localtime /etc/localtime

COPY ./main /bin/kk-direct

RUN chmod +x /bin/kk-direct

COPY ./config /config

COPY ./app.ini /app.ini

COPY ./lib/lua /lib/lua

COPY ./web /web

COPY ./static /static

COPY ./view /view

ENV LUA_PATH /lib/lua/?.lua;;

ENV KK_ENV_CONFIG /config/env.ini

VOLUME /config

EXPOSE 80

CMD kk-direct $KK_ENV_CONFIG

