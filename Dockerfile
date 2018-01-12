FROM golang:1.8
WORKDIR /Users/coop4/Mitchell/GoWebServerProject
RUN go get -v github.com/mitch-strong/GoWeb 
RUN go build -o main .
EXPOSE 8080
CMD ["/main"]
