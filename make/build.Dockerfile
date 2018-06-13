FROM golang:1.10-alpine

RUN apk add --no-cache git
RUN go get github.com/gorilla/sessions && \
	go get github.com/codegangsta/negroni && \
	go get github.com/joho/godotenv && \
	go get github.com/microcosm-cc/bluemonday && \
	go get gopkg.in/russross/blackfriday.v2 && \
	go get github.com/shurcooL/sanitized_anchor_name && \
	go get github.com/lib/pq && \
	go get github.com/gosimple/slug && \
	go get github.com/rainycape/unidecode && \
	go get github.com/oschwald/geoip2-golang && \
	go get github.com/tomasen/realip && \
	go get github.com/davecgh/go-spew/spew && \
	go get github.com/gorilla/mux && \
	go get golang.org/x/oauth2

RUN go get golang.org/x/crypto/ssh/terminal && \
	go get github.com/sirupsen/logrus