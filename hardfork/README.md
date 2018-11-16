hardfork
---
The repository aims to seperate the bitcoin-abc and bitcoin-sv after the timestamp `1542300000`(compared to mediantime). 

#### Notice:

This repository is under experiment state, so you should not test on your real coins unless you know what you do. **Important, I will not undertake the responsibility if you lost money.**

#### Usage:

1. Compile
    ```
    cd bitcoin-abc
    go build

    cd ../bitcoin-sv
    go build
    ```

2. Go ahead
    ```
    cd bitcoin-abc
    ./bitcoin-abc --privkey=******* --rpchost=127.0.0.1:8332 --rpcuser=rpc-user --rpcpassword=rpc-password --wait=600

    cd ../bitcoin-sv
    ./bitcoin-sv --privkey=******* --rpchost=127.0.0.1:8332 --rpcuser=rpc-user --rpcpassword=rpc-password --wait=600
    ```

