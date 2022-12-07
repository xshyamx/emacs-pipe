#!/bin/sh -e

i=0
n=10
s=3

while [ $i -lt $n ]; do
    i=$(expr $i + 1)
    printf "%d. %s\n" $i $RANDOM
    sleep 1
done
