#!/bin/bash

option=$1
delay=$2
jitters=$3
loss=$4

set -x
sudo tc qdisc del dev ens3 root

# option 1 means adding a default delay for everything
if [ $option -eq 1 ]; then
    echo "default delay mode, adding delay of ${delay}ms with ${jitters}ms jitters"
    sudo tc qdisc add dev ens3 root netem delay ${delay}ms ${jitters}ms 
fi

# option 2 adds a delay with loss and jitters
if [ $option -eq 2 ]; then
    echo "loss delay mode, adding delay of ${delay}ms with ${jitters}ms jitters, and ${loss}% loss"
    sudo tc qdisc add dev ens3 root netem delay ${delay}ms ${jitters}ms loss ${loss}%
fi


# instructions for adding a delay to a specific ip:port
# sudo tc qdisc add dev ens3 root handle 1: prio priomap 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2# sub the fifo qdisc with the prio one, with 3 bands by default, and send all traffic by default to band 1:3
# sudo tc qdisc add dev ens3 parent 1:1 handle 10: netem delay 100ms 10ms # add delay to band 0, does not affect anything yet as all traffic goes to band 2
# sudo tc filter add dev ens3 protocol ip parent 1:0 prio 1 u32 match ip dst 192.168.41.39/32 match ip dport 10000 0xffff flowid 1:1 # filter, put the chosen ip:port into band 0
