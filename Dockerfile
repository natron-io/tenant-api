# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /go/src/github.com/natron-io/tenant-api
RUN go get -d -v golang.org/x/net/html  
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tenant-api .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/natron-io/tenant-api/tenant-api ./
COPY --from=0 /go/src/github.com/natron-io/tenant-api/public ./public
EXPOSE 8000
CMD ["./tenant-api"]