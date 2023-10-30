FROM alpine:edge

RUN apk --no-cache add dnsmasq iproute2 iptables bash \
    && echo "conf-dir=/etc/dnsmasq.d,*.conf" > /etc/dnsmasq.conf

COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]