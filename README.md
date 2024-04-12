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

The result is shown in the graph below. TPS stands for transaction per second, BPS stands for bytes per second, and latency stands for the latency for reaching consesus.
| Test Case | broadcast TPS | gossip TPS | broadcast BPS | gossip BPS | broadcast latency | gossip latency |
| ------ | -------- | -------- | -------- | -------- | -------- |
| 1 | 879     | 773     | 450,045     | 396,031     | 998 | 1,112 |
| 2 | 数据     | 数据     | 数据     | 数据     |
| 3 | 数据     | 数据     | 数据     | 数据     |
| 4 | 数据     | 数据     | 数据     | 数据     |






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
