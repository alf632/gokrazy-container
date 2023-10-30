#!/bin/bash

MAC='00:c0:ca:b2:bd:9c' #'d8:3a:dd:0c:52:8f'
INTERFACE=$(ip a | grep -B1 "$MAC" | head -n1 | cut -f2 -d " " | sed -e "s/://")

function startClient(){
    sleep 5
    echo "starting Client"
    iwctl station $INTERFACE scan
    startDHCP&
}

function startDHCP() {
    while true
    do
        echo "starting dhclient on interface $INTERFACE"
        dhclient -4 -d -v $INTERFACE
        sleep 5
    done
}

service dbus start
startClient&
/usr/libexec/iwd
