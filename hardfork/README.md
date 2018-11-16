hardfork
---
这个仓库的目的是为了在区块的中位数时间戳大于等于`1542300000`时，有效分离bitcoin-abc和bitcoin-sv两条链。

#### 提醒:

这个仓库正处于活跃的开发中，如果您不是专业的开发人员，请不要将您真实的币用于该脚本中。

#### 使用说明:

1. 安装依赖：

    - 安装bitcoin-abc Linux系统下的客户端，安装文档请移步至: [https://github.com/Bitcoin-ABC/bitcoin-abc/tree/master/doc](https://github.com/Bitcoin-ABC/bitcoin-abc/tree/master/doc)；

    - 安装bitcoin-sv Linux系统下的客户端，安装文档请移步至：[https://github.com/bitcoin-sv/bitcoin-sv/tree/master/doc](https://github.com/bitcoin-sv/bitcoin-sv/tree/master/doc);

    - 安装Go语言环境：[Go](http://golang.org/) 1.8版本或者更新；

    - 安装Go包依赖：

        ```
        go get github.com/bcext/gcash
        go get github.com/bcext/cashutil
        go get github.com/qshuai/tcolor
        go get github.com/shopspring/decimal
        ```

2. 编译该工具: 该脚本在运行前，需要具有一些UTXO列表，当前的实现是通过硬编码一个UTXO列表在程序里面(主要是考虑在命令行手动输入大量utxo列表容易出错的缘故，故选择硬编码的方式)，那就是说每次在重新启动脚本的是否，您需要检查该UTXO列表是否可用。

    - 编译bitcoin-abc目录下的脚本程序，该脚本程序负责创建bitcoin-abc区块链上的特有交易(包含`OP_CHECKDATASIG`操作码)。

        ```
        cd $GOPATH/github.com/copernet/tools/hardfork/bitcoin-abc
        go install
        ```

    - 编译bitcoin-sv目录下的脚本程序，该脚本程序负责创建bitcoin-sv区块链上的特有交易(包含`OP_MUL`操作码)。

      ```
      cd $GOPATH/github.com/copernet/tools/hardfork/bitcoin-sv
      go install
      ```

    > 如果您没有没有安装Go语言环境和Go包依赖，我们为您准备了相应平台的二进制可执行文件，您可以直接运行，而不用经过上面上面繁琐的安装过程。二进制程序覆盖了linux和MacOS操作系统，您可以在bitcoin-abc和bitcoin-sv目录下找到它们，linux二进制包以`-linux`结尾，MaxOS二进制包以`osx`结尾。

3. 目前您需要检查bitcoin-abc和bitcoin-sv节点是否配置了`rpcuser`和`rpcpassword`选项，如果您执行的脚本程序的不在该节点上，您需要额外配置`rpcallowip`选项。在配置完相关选项之后，分别启动客户端bitcoin-abc和bitcoin-sv(强烈建议分开部署)。 在确定两个节点的数据已经同步到最高区块高度后(一般需要几小时到几十个小时，取决于您的网络环境和电脑配置)。准备工作已经完成，下面就可以执行我们的脚本程序了。

    下面是一个节点配置文件的示例(目前bitcoin-abc和bitcoin-sv的配置项没有什么差别)， 仅供参考:

    ```
    rpcuser=0XKRwRnyiY7NN2CfArU=
    rpcpassword=OSAqQIIp6XaPGuH53NvDlVPALQjaWksF4GPJyUimASq9
    txindex=1
    daemon=1
    rpcallowip=106.39.30.175
    ```

4. 启动脚本程序

    程序中几个选项的说明(具体的说明可以通过 `-h` 参看):

    - privkey: 钱包私钥，和硬编码的utxo需要对应，也就是这个私钥能够花费这些utxo；
    - rpchost: rpc server的ip和端口，默认为`127.0.0.1:8332`；
    - rpcuser: rpc server的用户名，用于用户rpc调用的认证；
    - rpcpassword：rpc server的密码，用户用户rpc调用的认证；
    - wait： 表示创建两笔交易之间的间隔时间，单位为秒(s)。需要注意的是如果您的utxo列表过少，同时wait参数过短，将会在节点内存池中产生许多关联交易，如果关联的深度大于25，将会在广播交易的时候失败。建以采用默认配置。

    ```
    // 创建bitcoin-abc独有共识的交易
    ./bitcoin-abc --privkey=******* --rpchost=127.0.0.1:8332 --rpcuser=rpc-user --rpcpassword=rpc-password --wait=600
    
    // 创建bitcoin-sv独有共识的交易
    ./bitcoin-sv --privkey=******* --rpchost=127.0.0.1:8332 --rpcuser=rpc-user --rpcpassword=rpc-password --wait=600
    ```

    程序运行的结果会直接在终端输出， 结果类似于以下的输出:

    ```
    available utxo: 1
    first transaction:
    	hash: 09f2e6c95f0fa31f357af90611803e2b35f9cc8256937db21298f7dc7cc8ad28
    rawtx: 0100000001333fb064838ad7685aee6a8650960c22be6ba59f0bc378b33da99ec889c0f79b010000006b483045022100aea45ac5e37125bb8442e634a8b128aa62636ccb9424a007b1fb897827f34ada02200bf1f4114c69336b1d44b76ea83511c0050c83458ecb861e1f3577176d6fe2c141210315eadd14931b1b41c903a291d63fd6aa50ba9523b5558db2857e18fd6df6a208ffffffff02220200000000000017a914fa8311a7df24c87fb759f3db2d2d763f0af5b071879c5c9800000000001976a9142b494c0b23d7aba3a01469fe265599d35179f8a388ac00000000
    second transction:
    	hash: 5b4ec490e265e908ef229f1e0e0caa94a251d105dac6be56d6ab749c9baa9e42
    	rawtx: 0100000001333fb064838ad7685aee6a8650960c22be6ba59f0bc378b33da99ec889c0f79b010000006b483045022100aea45ac5e37125bb8442e634a8b128aa62636ccb9424a007b1fb897827f34ada02200bf1f4114c69336b1d44b76ea83511c0050c83458ecb861e1f3577176d6fe2c141210315eadd14931b1b41c903a291d63fd6aa50ba9523b5558db2857e18fd6df6a208ffffffff02220200000000000017a914fa8311a7df24c87fb759f3db2d2d763f0af5b071879c5c9800000000001976a9142b494c0b23d7aba3a01469fe265599d35179f8a388ac00000000
    ```
