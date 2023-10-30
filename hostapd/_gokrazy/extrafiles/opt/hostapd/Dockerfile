FROM alpine:edge

RUN apk --no-cache add hostapd=2.10-r6 bash iproute2

COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
