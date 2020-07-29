FROM golang:alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN ls
RUN go build src/main.go
EXPOSE 80
CMD ["/app/main"]