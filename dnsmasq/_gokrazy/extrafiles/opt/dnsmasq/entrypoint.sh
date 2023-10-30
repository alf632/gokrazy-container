#!/bin/sh
MAC='d8:3a:dd:0c:52:8f'
INTERFACE=$(ip a | grep -B1 "$MAC" | head -n1 | cut -f2 -d " " | sed -e "s/://")

intNum=$(echo -n $INTERFACE | sed -e 's/wlan//')
tmpNum=$(($intNum - 1))
otherInterface=wlan${tmpNum#-}

while ! ip a show dev $INTERFACE | grep UP; do echo "waiting for hostapd on $INTERFACE"; sleep 5; done
ip a a 192.168.250.1/24 dev $INTERFACE

sysctl -w net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -s 192.168.250.0/24 -j MASQUERADE
iptables -P FORWARD ACCEPT
iptables -F FORWARD

echo "starting dnsmasq on interface $INTERFACE"
dnsmasq -d -i $INTERFACE --no-dhcp-interface=eth0 --no-dhcp-interface=$otherInterface --no-resolv