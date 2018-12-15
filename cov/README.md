#### Convert bitcoin cash address format between base58 and base32, supporting mainnet / testnet / regtest

##### How to use: 

1. Install Depedencies:

   ```
   go get github.com/bcext/cashutil
   go get github.com/bcext/gcash
   ```

2. clone this repository:

    ```
    cd $GOPATH
    git clone https://github.com/qshuai/Tools.git
    ```

3. install this tool:

    ```
    cd Tools/cov
    go install
    ```

4. Usage:

    ```
    cov qr35ennsep3hxfe7lnz5ee7j5jgmkjswssk2puzvgv mainnet
    // output: 1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX

    cov qqm2p7aglxw3dn7zzrdwmhd8lm2veypvlqu4yc6dtv testnet
    // output: mkVodiCHZ9AkZJLjSxYMvtweiYKr1wYP1w

    cov mfhy8vnokuzwygtBVEzPZxDDQj6bJJS45W regtest
    // output: bchreg:qqpp23mvg0aw5xt8lvql77yz9uuepyeqvsvcg2zqgv
    ```

##### ToDo:

- [ ] to support bitpay address format
