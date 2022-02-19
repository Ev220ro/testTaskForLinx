FROM golang:latest
MAINTAINER Rodionov Evgeniy
WORKDIR /app 
COPY . /app
CMD ["go", "run", "etl.go"]

