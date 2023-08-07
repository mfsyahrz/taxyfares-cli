FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o taxyfare .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/taxyfare .

CMD ["./taxyfare"]