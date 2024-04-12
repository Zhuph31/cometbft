# Comparison between broadcast and gossip
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

### Demo Video
This is the link to the demo video:

## Experiment
### Experiment Setting
We used 5 machines in total for the experiment. Geodec is deployed on the main machine, and four nodes will be started on each other machine. The machines are under the same subnet, and network conditions will be emulated through the network emulator module of Linux. Network emulators allows us to add artificial delays, jitters, and packet loss between the nodes.

### Geodec
We use geodec to perform the benchmark. With proper configuration, Geodec will log into the four machines, deploy cometbft, start a node on each machine, and keep sending transactions to the node during the benchmark period. We added new configurations for our four machines and set Geodec to use them. We also linked our version of cometbft with it so our modified version of cometbft will be deployed and used during the benchmark. The changes on Geodec are available at https://github.com/Zhuph31/geodec.

### Results
We conducted multiple tests, which covered the overall performance, network overhead, and network tolerance. 
### Overall Performance
For overall performance, we used data including throughput and latency. We emulated four kinds of network conditions.
1. Four nodes with 100ms delay and 10ms jitter between each pair of them
2. Four nodes with 200ms delay and 20ms jitter between each pair of them
3. Three nodes closely grouped together, the forth nodes is 300 ms delay and 30ms jitter away from the group
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
1. Broadcast is better than gossip at test cast 1, 2 and 4. Broadcast has higher throughput and takes less time to reach consensus, showing that broadcast can handle network delay and jitters better.
2. In test case 3, broadcast and gossip show almost the same performance. This is reasonable since under this network topology, the performance is contraint by the farthest node, therefore broadcast and gossip takes about the same time to reach this node and generates similar performance.

### Network Overhead
For network overhead, we conduct the benchmark without any artificial delay, counting the number of transactions and the network usage of cometbft nodes throughout the benchmark process. The result is availabe at the table below.
|          | broadcast | gossip    |
|----------|-----------|-----------|
| TPS      | 1,135 tx/s| 1,023 tx/s|
| Send     | 269 MB    | 338 MB    |
| Receive  | 219 MB    | 343 MB    |

From the table, we can see that broadcast achieved a higher TPS than gossip under the same network condition. What's more, the number of bytes sent and received for broadcast is much lower, meaning that broadcast mechanism uses less network resource to achieve better performance.

### Network Tolerance
For network tolerance, we used network emulator to add packet loss and monitor the performane under different level of packet losses. The resuls is in the table below.
| Packet Loss Rate | Broadcast TPS | Gossip TPS | Broadcast BPS | Gossip BPS | Broadcast Latency | Gossip Latency |
|------------------|---------------|------------|----------------|------------|-------------------|----------------|
| 5%               | 229 tx/s      | 443 tx/s   | 117,087 B/s    | 227,065 B/s| 1,819 ms          | 1,759 ms       |
| 10%              | 78 tx/s       | 0 tx/s     | 39,694 B/s     | 0 B/s      | 2,197 ms          | -              |
| 20%              | 78 tx/s       | 0 tx/s     | 39,694 B/s     | 0 B/s      | 2,197 ms          | -              |

From the table, we can tell that:
1. Gossip has a better performance when handling low packet loss rate like 5%.
2. When packet loss rate goes up, gossip cannot reach consensus in time and does not commit any blocks, while broadcast can still commit some blocks.

### Usage
Since only the internal mechanism is changed, the uasge of this repo is the same as the original cometbft repo. You may find what you need to do to have the code running here https://github.com/cometbft/cometbft.
For testing, please refere to https://github.com/Zhuph31/geodec.

### 
6. How to use?
    1. Instructions for a local setup to host the application
    2. Required libraries for running it - provide installation instructions, links are acceptable. Scripts are preferred wherever possible. 
    3. How to run it? - walk through the different features users can use.
    4. If there are different types of users, specify what each can do.
    5. If Metamask is used, indicate which testnet to link to.
7. How to contribute?
    1. Architecture
    2. Local setup instructions
        1. Include required dependencies installation instructions, if any.
        2. Setup for testing, how to run tests.
        3. If a smart contract, how to deploy a new contract?
8. Include an appropriate license for the project.
9. If this repository used components from other projects or is a fork:
    1. Give credits to upstream repositories 
    2. Clearly specify what is built on top of them.
- Answer all relevant questions, adding more details for completeness as needed.
- In addition to the README, we expect code to be present in the same repository.
    - The code has to be well commented and easy to follow.
    - Please use capabilities of GitHub - clearly defined PRs, commits and reviews on PRs from your team members.
    - We expect individual contributions in terms of code from everyone in team and this data might be used for re-weighting the individual grades.
