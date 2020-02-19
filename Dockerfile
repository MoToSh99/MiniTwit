# docker run --name webserver -p 5000:5000 -p 5001:5001  minitwit
FROM golang:latest

RUN go get github.com/gorilla/mux
RUN go get github.com/Jeffail/gabs
RUN go get github.com/gorilla/securecookie
RUN go get github.com/jinzhu/gorm
RUN go get github.com/jinzhu/gorm/dialects/sqlite
RUN go get github.com/jinzhu/inflection
RUN go get golang.org/x/crypto/bcrypt
RUn go get github.com/jinzhu/gorm/dialects/mssql



WORKDIR /src

COPY . .

EXPOSE 5000
EXPOSE 5001

RUN go build /src/server.go

CMD ["/src/server"]