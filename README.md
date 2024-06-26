# Broadcast vs. Gossip: A Comparative Analysis in CometBFT Consensus"

## Context
Byzantine fault tolerance (BFT) enables distributed systems to maintain functionality even if some components fail or act maliciously. It ensures system integrity and consensus by replicating data across multiple nodes and employing algorithms to detect and mitigate faulty behavior.

Our objective was to create a comprehensive analysis of gossip and broadcast communication protocols by using both on a predetermined Byzantine fault tolerance (BFT) algorithm which would be CometBFT in this case. By delving into their respective designs, consensus mechanisms, and other various metrics, the project seeks to compare their trade-offs. Which we have completed successfully.

## Problem Statement
The problem we are currently facing is the lack of a unified platform along with the metrics that can be used to compare both gossip and broadcast communication protocols. By comparing the differences between gossip and broadcast communication protocols within the context of Byzantine Fault Tolerances (BFT), we can determine their performance and metrics. Addressing this problem will allow us to understand the strengths and weaknesses of gossip and broadcast protocols in BFT systems which will, in turn, enable more informed decision-making when selecting protocols for achieving fault tolerance. We can also use this information to improve resilience and efficiency in distributed systems such as blockchain systems.

## Solution
To compare and benchmark gossip and broadcast in a fair and meaningful way, we plan to put these two methods into the same implementation and compare them. For this, we choose to use CometBFT as a base. CometBFT can be easily deployed locally and can form a network of multiple nodes on the same machine with the help of docker.
CometBFT adopts the gossip protocol by default, and we aim to implement a broadcast protocol on top of it. After the broadcast protocol is implemented, we can compare and benchmark these two methods using the modified CometBFT and this can guarantee that the experiments for each protocol will be under the same circumstances. Overall, the key features of the modified CometBFT consist of the complete implementation of a broadcast communication protocol as opposed to the originally implemented gossip communication protocol.
For benchmark, we employed Geodec. Geodec is a repository designed to investigate the influence of geographic locations on consensus algorithms. In our case, we simply use Geodec to test the performance of the gossip and broadcast mechanisms.

### Code Implementation
#### Vote Broadcast

The updated reactor.go into the internal/concensus directory manages the broadcasting and synchronization of data and voting for CometBFT. The newly created functions BroadcastDataRoutine and BroadcastDataForCatchup handle the dissemination of block parts, proposals, and catch-up data to ensure network consistency. The newly created BroadcastVotesRoutine handles the broadcasting of votes based on consensus steps and heights, employing a sleep mechanism when no votes are available. Newly created broadcastVotesForHeight assists in selecting and sending specific votes pertinent to the current consensus state. The modified queryMaj23Routine focuses on querying and broadcasting two-thirds majority votes for prevotes, precommits, and proposal POLs, contributing to the network's consensus integrity and progress. This essentially sums up the modified changes and the responsibility of the Broadcast voting that has been implemented.

#### Transaction Broadcast
By default, the node that receives the transaction would send the transaction to all its peers, and peers would keep on gossiping about the transaction. To keep the other nodes from relaying the transaction, we record the sender when receiving a transaction and compare it against all the peers. If the transaction is sent by one of the peers, we believe that it is the sender's responsibility to broadcast the transaction. Otherwise, we assume it is sent by a client and start broadcasting the transaction.


## Packages to install
Please refer to the official cometbft guide for install: https://docs.cometbft.com/v0.38/guides/install


## Setup
Please refer to the official cometbft guide for quick start: https://docs.cometbft.com/v0.38/guides/quick-start

This is a quick start guide. If you have a vague idea about how CometBFT works and want to get started right away, continue.

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
For network overhead, we conduct the benchmark without any artificial delay, counting the number of transactions and the network usage of cometbft nodes throughout the benchmark process. The result is available in the table and graphs below.
|          | broadcast | gossip    |
|----------|-----------|-----------|
| TPS      | 1,135 tx/s| 1,023 tx/s|
| Send     | 269 MB    | 338 MB    |
| Receive  | 219 MB    | 343 MB    |

![image](https://github.com/Zhuph31/cometbft/assets/50798194/bd2e6b76-76b1-405c-b26a-cfc3058a4994)
![image](https://github.com/Zhuph31/cometbft/assets/50798194/083392f8-3e13-439d-85d0-0463955064a7)
![image](https://github.com/Zhuph31/cometbft/assets/50798194/9c57a892-f3f8-4dd7-b9db-27627d0d983b)

From the result, we can see that broadcast achieved a higher TPS than gossip under the same network condition. What's more, the number of bytes sent and received for broadcast is much lower, meaning that the broadcast mechanism uses fewer network resources to achieve better performance.

### Network Tolerance
For network tolerance, we used a network emulator to add packet loss and monitor the performance under different levels of packet loss. The results are in the table and graphs below.
| Packet Loss Rate | Broadcast TPS | Gossip TPS | Broadcast BPS | Gossip BPS | Broadcast Latency | Gossip Latency |
|------------------|---------------|------------|----------------|------------|-------------------|----------------|
| 5%               | 229 tx/s      | 443 tx/s   | 117,087 B/s    | 227,065 B/s| 1,819 ms          | 1,759 ms       |
| 10%              | 78 tx/s       | 0 tx/s     | 39,694 B/s     | 0 B/s      | 2,197 ms          | -              |
| 20%              | 78 tx/s       | 0 tx/s     | 39,694 B/s     | 0 B/s      | 2,197 ms          | -              |

![image](https://github.com/Zhuph31/cometbft/assets/50798194/fb6f119a-a6fa-4954-9569-c41af5b40b65)
![image](https://github.com/Zhuph31/cometbft/assets/50798194/3895764a-2acf-4c7b-92d3-5aed47ccd18c)

From the result, we can tell that:
1. Gossip has a better performance when handling a low packet loss rate like 5%.
2. When packet loss rate goes up, gossip cannot reach consensus in time and does not commit any blocks, while broadcast can still commit some blocks.

### Usage
Since only the internal mechanism is changed, the usage of this repo is the same as the original cometbft repo. For building, installing, and local testing, please refer to  https://github.com/cometbft/cometbft.
For benchmarking, please refer to https://github.com/Zhuph31/geodec main branch for instructions.

### Credits
This repo is forked from https://github.com/cometbft/cometbft.
The modification on top of the original repo is the vote broadcast and transaction broadcast specified in previous sections.
### References

References
1. Cometbft. (Feb 25, 2024). CometBFT. GitHub. https://github.com/cometbft/cometbft
2. Cometbft docs. (Feb 25, 2024). CometBFT official documents.
https://docs.cometbft.com/v0.38/
3. Broadcast gossip algorithms for consensus. (2009, July 1). IEEE Journals & Magazine |
IEEE Xplore. https://ieeexplore.ieee.org/abstract/document/4787122
4. Netem (Network Emulator). [OpenWrt Wiki]. (n.d.). https://openwrt.org/docs/guide-user/network/traffic-shaping/sch_netem 
5.GeoDecConsensus. (n.d.). Geodecconsensus/GEODEC. GitHub. https://github.com/GeoDecConsensus/geodec 

