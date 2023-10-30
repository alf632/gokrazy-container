#!/bin/bash

MAC='d8:3a:dd:0c:52:8f'
INTERFACE=$(ip a | grep -B1 "$MAC" | head -n1 | cut -f2 -d " " | sed -e "s/://")

echo "starting hostapd on $INTERFACE"
hostapd -i $INTERFACE /config/hostapd.conf 