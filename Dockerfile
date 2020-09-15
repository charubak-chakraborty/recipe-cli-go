FROM golang:latest

LABEL maintainer="Charubak Chakraborty <charubak.chakraborty@gmail.com>"

RUN mkdir /app

COPY . /app

ARG filename

# COPY fixture
COPY ${filename} /app/temp/input.json

WORKDIR /app

RUN go get -d -v ./...

RUN go build 

CMD ./app