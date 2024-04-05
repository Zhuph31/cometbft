#!/bin/bash

option=$1
delay=$2
jitters=$3

sudo tc qdisc del dev ens3 root

# option 1 means adding a default delay for everything
if [ $option -eq 1 ]; then
    echo "default delay mode, adding delay of ${delay}ms with ${jitters}ms jitters"
    sudo tc qdisc add dev ens3 root netem delay ${delay}ms ${jitters}ms 
fi
