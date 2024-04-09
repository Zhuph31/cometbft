# Comparison between broadcast and gossip
1. Title of project 
2. Context
3. Problem Statement
## Context
## Problem Statement
## Solution
### Vote Broadcast
### Transaction Broadcast
By default, the node that receives the transaction would send the transaction to all its peers, and peers would keep on gossiping about the transaction. To keep the other nodes from relaying the transaction, we record the sender when receiving a transaction and compare it against all the peers. If the transaction is sent by one of the peers, we believe that it is the sender's responsibility to broadcast the transaction. Otherwise, we assume it is sent by a client and start broadcasting the transaction.

### Demo Video
This is the link to the demo video:

## Experiment
### Experiment Setting
We used 5 machines in total for the experiment. Geodec is deployed on the main machine, and four nodes will be started on each other machine. The machines are under the same subnet, and network conditions will be emulated through the network emulator module of Linux.

### Geodec
We use geodec to perform the benchmark. With proper configuration, geodec will log into the four machines, deploy cometbft, start a node on each machine, and keep sending transactions to the node during the benchmark period. We added new configurations for our four machines and set geodec to use them. We also linked our version of cometbft with it so our modified version of cometbft will be deployed and used during the benchmark. The changes on geodec is available at https://github.com/Zhuph31/geodec.


    5. Include graphs from benchmarking.
5. How to use?
    1. Instructions for a local setup to host the application
    2. Required libraries for running it - provide installation instructions, links are acceptable. Scripts are preferred wherever possible. 
    3. How to run it? - walk through the different features users can use.
    4. If there are different types of users, specify what each can do.
    5. If Metamask is used, indicate which testnet to link to.
6. How to contribute?
    1. Architecture
    2. Local setup instructions
        1. Include required dependencies installation instructions, if any.
        2. Setup for testing, how to run tests.
        3. If a smart contract, how to deploy a new contract?
7. Include an appropriate license for the project.
8. If this repository used components from other projects or is a fork:
    1. Give credits to upstream repositories 
    2. Clearly specify what is built on top of them.
- Answer all relevant questions, adding more details for completeness as needed.
- In addition to the README, we expect code to be present in the same repository.
    - The code has to be well commented and easy to follow.
    - Please use capabilities of GitHub - clearly defined PRs, commits and reviews on PRs from your team members.
    - We expect individual contributions in terms of code from everyone in team and this data might be used for re-weighting the individual grades.
