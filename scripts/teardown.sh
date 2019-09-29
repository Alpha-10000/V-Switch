#!/bin/sh

if [ $# -lt 1 ]
then
    echo "Usage: teardown.sh NAME"
else
    for name in "$@"
    do
	echo "Destroying namespace $name..."
	ip netns del $name
    done
fi
