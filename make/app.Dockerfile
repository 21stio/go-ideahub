FROM golang:1.10-alpine

RUN apk add --no-cache postgresql-client
COPY GeoLite2-City.mmdb GeoLite2-City.mmdb
COPY routes routes
COPY public public
COPY main main
COPY make make
COPY Makefile Makefile
RUN apk add --no-cache make

CMD ./main