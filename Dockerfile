FROM golang:1.14.1-alpine AS builder
RUN apk --no-cache add build-base git gcc
COPY . /go/src/app
WORKDIR /go/src/app/
RUN go mod download
RUN cd /go/src/app && CGO_ENABLE=0 GOARCH=amd64 go build -a -o main .

# final stage
FROM alpine
WORKDIR /app/
COPY --from=builder /go/src/app/main .
ENTRYPOINT ["./main"]