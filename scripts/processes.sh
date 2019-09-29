#!/bin/sh

if [ $# -lt 1 ]
then
    echo "Specify namespace ID"
else
    ps -eo pid,user,netns,args --sort user | grep $1
fi
