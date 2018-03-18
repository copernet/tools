## big size block(BCH) testing

This repository is under experiment. The current problem is that the client will not pack the expected block size even if the mempool size has up to 190+ megabyte.

### Env

- Ubuntu 16.04.3 LTS
- Bitcoin-ABC v0.16.2

> startup command: ./bitcoind -excessiveblocksize=64000000 -blockmaxsize=32000000 -testnet

