FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o crutchdb cmd/server/main.go

WORKDIR /app

RUN cp /build/crutchdb crutchdb

EXPOSE 2357

CMD ["/app/crutchdb"]