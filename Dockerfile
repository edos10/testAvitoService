FROM golang:latest as builder
COPY go.mod go.sum /go/src/github.com/edos10/test_avito_service/
WORKDIR /go/src/github.com/edos10/test_avito_service/
RUN go mod download
COPY . /go/src/github.com/edos10/test_avito_service/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/test_avito_service github.com/edos10/test_avito_service/

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/edos10/test_avito_service/build/test_avito_service /usr/bin/test_avito_service
ENTRYPOINT ["/usr/bin/test_avito_service"]