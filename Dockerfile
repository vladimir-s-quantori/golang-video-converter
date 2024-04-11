FROM golang:latest as builder
LABEL authors="VladimirSemenovich"

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
#RUN go env -w GO111MODULE=on

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /vidConv

EXPOSE 8080

CMD ["/vidConv"]

