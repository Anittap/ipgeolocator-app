FROM golang:alpine as builder

WORKDIR /app

COPY main.go .

RUN go mod init crm-app
RUN go get github.com/bradfitz/gomemcache/memcache
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/gin-gonic/gin
RUN go get github.com/aws/aws-sdk-go-v2
RUN go get github.com/aws/aws-sdk-go-v2/config
RUN go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
RUN go build -o main .

FROM alpine:latest 
ENV USER=gouser
RUN mkdir /goapp
RUN adduser -D -h /goapp -s /bin/sh $USER
WORKDIR /goapp

RUN apk update && apk add --no-cache \
    ca-certificates \
    bash \
    && rm -rf /var/cache/apk/*

COPY --from=builder /app/main .

RUN chown -R $USER:$USER . 

EXPOSE 8080
USER $USER

ENTRYPOINT ["./main"]
