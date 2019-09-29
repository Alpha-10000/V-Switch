#!/bin/sh

if [ $# -lt 3 ]
then
    echo "Usage: add-intf.sh NAME NS1 NS2"
else
    vethA="$1-1-$2"
    vethB="$1-2-$3"
    echo "Creating interface '$vethA' in namespace '$2'..."
    echo "Creating interface '$vethB' in namespace '$3'..."
    ip link add "$vethA" netns "$2" type veth peer name "$vethB" netns "$3"
fi
