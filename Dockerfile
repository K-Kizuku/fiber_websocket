FROM golang:1.20.1

ARG APP_NAME

RUN mkdir -p /go/src/${APP_NAME}/
COPY . /go/src/${APP_NAME}/
WORKDIR /go/src/${APP_NAME}/

RUN go install

RUN go build -o /go/bin/air github.com/cosmtrek/air

CMD air