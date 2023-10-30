from debian:bullseye-slim

RUN groupadd netdev
RUN apt update && apt install -y dbus iwd rfkill isc-dhcp-client

COPY entrypoint.sh /entrypoint.sh
COPY iwd.conf /etc/iwd/main.conf

cmd /entrypoint.sh
