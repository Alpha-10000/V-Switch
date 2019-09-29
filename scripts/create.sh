#!/bin/sh

if [ $# -lt 1 ]
then
    echo "Usage: create.sh NAME"
else
    for name in "$@"
    do
	echo "Creating namespace $name..."
	ip netns add $name
    done
fi
