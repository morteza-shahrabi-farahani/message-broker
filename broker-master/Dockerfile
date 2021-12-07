FROM golang

WORKDIR /

COPY . .
RUN go build -o server server.go

CMD ["/server/server"]


