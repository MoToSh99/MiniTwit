# docker run --name webserver -p 5000:5000 -p 5001:5001  minitwit
FROM golang:latest

RUN go get github.com/gorilla/mux
RUN go get github.com/Jeffail/gabs
RUN go get github.com/gorilla/securecookie
RUN go get github.com/jinzhu/gorm
RUN go get github.com/jinzhu/gorm/dialects/sqlite
RUN go get github.com/jinzhu/inflection
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/jinzhu/gorm/dialects/mssql
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promauto
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/alecthomas/template
RUN go get github.com/swaggo/http-swagger
RUN go get github.com/bshuster-repo/logrus-logstash-hook
RUN go get github.com/sirupsen/logrus
RUN go get github.com/joho/godotenv
RUN go get github.com/lib/pq
RUN go get github.com/unrolled/secure
RUN go get github.com/jinzhu/gorm/dialects/postgres

WORKDIR /src

COPY /app .

EXPOSE 5000
EXPOSE 5001


RUN go build /src/main.go

CMD ["/src/main"]
