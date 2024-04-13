# Broadcast vs. Gossip: A Comparative Analysis in CometBFT Consensus"
1. Title of project 
2. Context
3. Problem Statement
## Context
## Problem Statement
## Solution
We start by understanding the code of cometbft, and learning how it works on mempool and consensus levels. Then we implemented the broadcast mechanism on top of cometbft, tested it, and ensured it was working in the correct way.
For benchmark, we employed Geodec. Geodec is a repository designed to investigate the influence of geographic locations on consensus algorithms. In our case, we simply use Geodec to test the performance of the gossip and broadcast mechanisms.
### Code Implementation
#### Vote Broadcast
#### Transaction Broadcast
By default, the node that receives the transaction would send the transaction to all its peers, and peers would keep on gossiping about the transaction. To keep the other nodes from relaying the transaction, we record the sender when receiving a transaction and compare it against all the peers. If the transaction is sent by one of the peers, we believe that it is the sender's responsibility to broadcast the transaction. Otherwise, we assume it is sent by a client and start broadcasting the transaction.

## Experiment
### Experiment Setting
We used 5 machines in total for the experiment. Geodec is deployed on the main machine, and four nodes will be started on each other machine. The machines are under the same subnet, and network conditions will be emulated through the network emulator module of Linux. Network emulators allow us to add artificial delays, jitters, and packet loss between the nodes.

### Geodec
We use geodec to perform the benchmark. With proper configuration, Geodec will log into the four machines, deploy cometbft, start a node on each machine, and keep sending transactions to the node during the benchmark period. We added new configurations for our four machines and set Geodec to use them. We also linked our version of cometbft with it so our modified version of cometbft will be deployed and used during the benchmark. The changes on Geodec are available at https://github.com/Zhuph31/geodec.

### Results
We conducted multiple tests, which covered the overall performance, network overhead, and network tolerance. 
### Overall Performance
For overall performance, we used data including throughput and latency. We emulated four kinds of network conditions.
1. Four nodes with 100ms delay and 10ms jitter between each pair of them
2. Four nodes with 200ms delay and 20ms jitter between each pair of them
3. Three nodes closely grouped together, the fourth node has a 300 ms delay and 30 ms jitter away from the group
4. Four nodes, distributed in pairs on each side, with a 100ms delay and 10ms jitter in between.

The result is shown in the table and bar graph below. TPS stands for transaction per second, BPS stands for bytes per second, and latency stands for the latency for reaching consesus.
| Test Case | broadcast TPS | gossip TPS | broadcast BPS | gossip BPS | broadcast latency | gossip latency |
| --------- | -------------- | ---------- | -------------- | ---------- | ----------------- | -------------- |
| 1         | 879            | 773        | 450,045        | 396,031    | 998               | 1,112          |
| 2         | 209            | 164        | 107,107        | 84,188     | 21,707            | 27,778         |
| 3         | 1,007          | 1,009      | 515,525        | 516,365    | 1,275             | 1,207          |
| 4         | 1,217          | 954        | 623,290        | 488,222    | 1,019             | 652            |

![image](https://github.com/Zhuph31/cometbft/assets/50798194/98b1401a-eb4f-4935-8ff4-21202c9fc33e)
![image](https://github.com/Zhuph31/cometbft/assets/50798194/0f2bf335-a402-4352-8d7a-af9567c53e1b)
![image](https://github.com/Zhuph31/cometbft/assets/50798194/b3aaf29e-3452-4669-9eeb-fb01e96a89bd)

From the results we can see that:
1. Broadcast is better than gossip in test cases 1, 2, and 4. Broadcast has higher throughput and takes less time to reach consensus, showing that broadcast can handle network delay and jitters better.
2. In test case 3, broadcast and gossip show almost the same performance. This is reasonable since under this network topology, the performance is constrained by the farthest node, therefore broadcast and gossip take about the same time to reach this node and generate similar performance.

### Network Overhead
For network overhead, we conduct the benchmark without any artificial delay, counting the number of transactions and the network usage of cometbft nodes throughout the benchmark process. The result is available at the table below.
|          | broadcast | gossip    |
|----------|-----------|-----------|
| TPS      | 1,135 tx/s| 1,023 tx/s|
| Send     | 269 MB    | 338 MB    |
| Receive  | 219 MB    | 343 MB    |

From the table, we can see that broadcast achieved a higher TPS than gossip under the same network condition. What's more, the number of bytes sent and received for broadcast is much lower, meaning that the broadcast mechanism uses fewer network resources to achieve better performance.

### Network Tolerance
For network tolerance, we used a network emulator to add packet loss and monitor the performance under different levels of packet loss. The results are in the table below.
| Packet Loss Rate | Broadcast TPS | Gossip TPS | Broadcast BPS | Gossip BPS | Broadcast Latency | Gossip Latency |
|------------------|---------------|------------|----------------|------------|-------------------|----------------|
| 5%               | 229 tx/s      | 443 tx/s   | 117,087 B/s    | 227,065 B/s| 1,819 ms          | 1,759 ms       |
| 10%              | 78 tx/s       | 0 tx/s     | 39,694 B/s     | 0 B/s      | 2,197 ms          | -              |
| 20%              | 78 tx/s       | 0 tx/s     | 39,694 B/s     | 0 B/s      | 2,197 ms          | -              |

From the table, we can tell that:
1. Gossip has a better performance when handling a low packet loss rate like 5%.
2. When packet loss rate goes up, gossip cannot reach consensus in time and does not commit any blocks, while broadcast can still commit some blocks.

### Usage
Since only the internal mechanism is changed, the usage of this repo is the same as the original cometbft repo. For building, installing, and local testing, please refer to  https://github.com/cometbft/cometbft.
For benchmarking, please refer to https://github.com/Zhuph31/geodec main branch for instructions.

### Credits
This repo is forked from https://github.com/cometbft/cometbft.
The modification on top of the original repo is the vote broadcast and transaction broadcast specified in previous sections.
