FROM alpine:3.7

COPY GeoLite2-City.mmdb GeoLite2-City.mmdb
COPY public public
COPY main main

CMD main