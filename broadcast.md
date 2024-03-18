# Possible Solution to swtich from gossip to broadcast

By default, cometbft nodes do peer exchange through PEX, and each node would have a number of connected nodes. Each node would relay TXs to connected nodes. (Gossip)

One possible way to simulate a broadcast is to let all nodes connect with each other, and disable the TXs relay function so the TXs will be broadcasted to the whole network by one node.

1. Cometbft allow you to specify persistent peers that would be connected all the time. 
https://github.com/Zhuph31/cometbft/blob/main/docs/explanation/core/using-cometbft.md#persistent-peer
Adding all other nodes as persistent peers to each node when starting up would allow each node to talk to all other nodes.

2. Cometbft has a boolean flag "broadcast" that controls whether the node would relay TXs to other connected peers. It seems like the flag cannot be controlled through configuration or commandline argument so you may have to modify it a little.
https://github.com/Zhuph31/cometbft/blob/main/config/config.go#L873

The remain problem is that I have not tested whether a node would send TXs to all connected peers on proposal. If so, setting all nodes as persistent peers and disabling the "broadcast" function should allow us to simulate a broadcast.


## Other find outs
1. When using broadcast, PEX can be turned off by controlling this boolean.
https://github.com/Zhuph31/cometbft/blob/main/config/config.go#L726

2. Metrics here can be used to design the experiment. Some of them like p2p_message_send_bytes_total, p2p_message_receive_bytes_total are pretty detailed.
https://github.com/Zhuph31/cometbft/blob/main/docs/explanation/core/metrics.md

3. Disable empty blocks. Cometbft creates emtpy blocks by default and it may be easier to monitor with empty blocks turned off.
https://github.com/Zhuph31/cometbft/blob/main/docs/explanation/core/using-cometbft.md#no-empty-blocks
