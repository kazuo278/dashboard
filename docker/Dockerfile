FROM golang:1.21-alpine as builder

WORKDIR /workspace

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download

COPY . /workspace/
RUN CGO_ENABLED=0 GOOS=linux go build -o dashboard

FROM alpine:3.17

RUN apk add tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone
COPY --from=builder /workspace/dashboard /app/dashboard
COPY static-content/ /static-content/
RUN chmod +x /app/dashboard
ENTRYPOINT [ "/app/dashboard" ]