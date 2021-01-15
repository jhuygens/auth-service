FROM golang:alpine as builder

RUN mkdir /build
ADD ./ /build 
WORKDIR /build
RUN env GOOS=linux GOARCH=386 go build -o main .

FROM alpine:latest

RUN mkdir -p /app && adduser -S -D -H -h /app appuser && chown -R appuser /app
COPY --from=builder /build/main /app/
USER appuser
EXPOSE 9091
WORKDIR /app
CMD ["./main"]