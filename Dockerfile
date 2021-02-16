FROM golang:1.13.5-alpine3.10 AS builder
LABEL maintainer="liangdas <1587790525@qq.com>" github="https://github.com/liangdas"
RUN apk --no-cache add git
WORKDIR /build

ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on
ENV GOSUMDB off
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o tetris-server .

LABEL version="1.0.0"

FROM alpine:3.10 AS final

WORKDIR /app
COPY --from=builder /build/tetris-server /app/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=nats:2.1.9-alpine3.12 /usr/local/bin/nats-server /nats-server
COPY --from=consul /bin/consul /consul

COPY start.sh /app/start.sh
COPY stop.sh /app/stop.sh
COPY docker-entrypoint.sh /docker-entrypoint.sh
COPY bin/conf /app/bin/conf
COPY static /app/static
COPY tetrisconfig /app/tetrisconfig

CMD source /docker-entrypoint.sh
CMD source /start.sh
CMD source /stop.sh
EXPOSE 6653 6563 6565
RUN chmod a+x /app/start.sh
RUN chmod a+x /app/stop.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]