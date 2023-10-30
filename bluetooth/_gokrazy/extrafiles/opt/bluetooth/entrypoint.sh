#!/bin/bash

dev=hci0

# reset bluetooth adapter by restarting it
hciconfig $dev down
hciconfig $dev up

# from https://github.com/RPi-Distro/pi-bluetooth/blob/master/usr/bin/bthelper
if ( /usr/bin/hcitool -i $dev dev | grep -q -E '\s43:4[35]:|AA:AA:AA:AA:AA:AA' ); then
    SERIAL=`cat /proc/device-tree/serial-number | cut -c9-`
    B1=`echo $SERIAL | cut -c3-4`
    B2=`echo $SERIAL | cut -c5-6`
    B3=`echo $SERIAL | cut -c7-8`
    BDADDR=`printf '0x%02x 0x%02x 0x%02x 0xeb 0x27 0xb8' $((0x$B3 ^ 0xaa)) $((0x$B2 ^ 0xaa)) $((0x$B1 ^ 0xaa))`

    /usr/bin/hcitool -i $dev cmd 0x3f 0x001 $BDADDR
    /usr/bin/hciconfig $dev reset
else
    echo Raspberry Pi BDADDR already set
fi

# Route SCO packets to the HCI interface (enables HFP/HSP)
/usr/bin/hcitool -i $dev cmd 0x3f 0x1c 0x01 0x02 0x00 0x01 0x01 > /dev/null

# start services
service dbus start
service bluetooth start

# wait for startup of services
msg="Waiting for services to start..."
time=0
echo -n $msg
while [[ "$(pidof start-stop-daemon)" != "" ]]; do
    sleep 1
    time=$((time + 1))
    echo -en "\r$msg $time s"
done
echo -e "\r$msg done! (in $time s)"

echo "starting bt-agent"
/usr/bin/bt-agent -c NoInputNoOutput -p /opt/bt-agent.conf
#/usr/bin/bluetoothctl power off; sleep 1; /usr/bin/bluetoothctl power on

while true
do
  sleep 1
done
