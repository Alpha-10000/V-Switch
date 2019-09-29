#!/bin/sh

if [ $# -lt 2 ]
then
    echo "Usage: del-intf.sh NAME NS"
else
    echo "Deleting interface $1 from namespace $2..."
    ip netns exec "$2" ip link set "$1" netns 1
    ip link del "$1"
fi
