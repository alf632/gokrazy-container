FROM debian:bullseye-slim as builder

#environment variables
ENV BLUEZ_VERSION 5.69

RUN apt-get update && apt-get install -y \
    build-essential wget systemd \
    libical-dev libdbus-1-dev libglib2.0-dev libreadline-dev libudev-dev udev

RUN wget -P /tmp/ https://www.kernel.org/pub/linux/bluetooth/bluez-${BLUEZ_VERSION}.tar.gz \
 && tar xf /tmp/bluez-${BLUEZ_VERSION}.tar.gz -C /tmp \
#compile bluez
 && cd /tmp/bluez-${BLUEZ_VERSION} \
 && ./configure --prefix=/usr \
    --mandir=/usr/share/man \
    --sysconfdir=/etc \
    --localstatedir=/var \
    --enable-library \
    --enable-experimental \
    --enable-maintainer-mode \
    --enable-deprecated \
    --enable-tools \
    --disable-manpages \
 && make \
#install bluez tools
 && make install


FROM debian:bullseye-slim

RUN apt update && apt install -y dbus bluez-tools

RUN useradd -m bluezuser \
 && passwd -d bluezuser

COPY --from=builder /usr/bin/bluetoothctl /usr/bin/btmon /usr/bin/rctest /usr/bin/l2test /usr/bin/l2ping \
                    /usr/bin/bluemoon /usr/bin/hex2hcd /usr/bin/mpris-proxy /usr/bin/btattach \
                    /usr/bin/hciattach /usr/bin/hciconfig /usr/bin/hcitool /usr/bin/hcidump /usr/bin/rfcomm \
                    /usr/bin/sdptool /usr/bin/ciptool /usr/bin/
COPY --from=builder /usr/share/dbus-1/system.d/bluetooth.conf /etc/dbus-1/system.d/
COPY --from=builder /usr/libexec/bluetooth/bluetoothd /usr/libexec/bluetooth/obexd /usr/libexec/bluetooth/
COPY --from=builder /usr/share/dbus-1/services/org.bluez.obex.service /usr/share/dbus-1/services/
COPY --from=builder /usr/share/dbus-1/system-services/org.bluez.service /usr/share/dbus-1/system-services/
COPY --from=builder /usr/include/bluetooth/* /usr/include/bluetooth/
COPY --from=builder /usr/lib/pkgconfig/bluez.pc /usr/lib/pkgconfig/
COPY --from=builder /usr/lib/bluetooth/plugins/external-dummy.so /usr/lib/bluetooth/plugins/
COPY --from=builder /usr/lib/bluetooth/plugins/external-dummy.la /usr/lib/bluetooth/plugins/
COPY --from=builder /lib/systemd/system/bluetooth.service /lib/systemd/system
COPY --from=builder /etc/bluetooth /etc/bluetooth

# setup startup script
COPY entrypoint.sh .
COPY bluezuser.conf /etc/dbus-1/system.d/
COPY bluetooth.init /etc/init.d/bluetooth
RUN chmod +x ./entrypoint.sh
CMD ./entrypoint.sh

