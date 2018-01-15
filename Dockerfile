FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./Web /go/src/github.com/mitch-strong/GoWeb/Web
WORKDIR /go/src/github.com/mitch-strong/GoWeb/Web

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
	
EXPOSE 8080