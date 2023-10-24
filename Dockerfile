FROM alpine:3.17

RUN apk add tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone
COPY dashboard /app/dashboard
COPY static-content/ /static-content
RUN chmod +x /app/dashboard
ENTRYPOINT [ "/app/dashboard" ]