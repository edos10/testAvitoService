FROM golang:latest AS builder

WORKDIR /app

COPY . ./

RUN go mod download && go build -o /server



#FROM alpine:3
#
#COPY --from=builder server /bin/server

ENTRYPOINT ["/server"]