## Mine Big Size Blocks(BCH) Testing

This repository is under experiment. The current problem is that kafka is not suitable for transfering a big size message(up to 90 megabyte). Up to now, I have mined several blocks with about 16 megabyte size.

### Env:

- Ubuntu 16.04.3 LTS
- Bitcoin-ABC v0.16.2

### Run bitcoind client:

```
cd path/to/bitcoin-abc/src
./bitcoind -testnet -relaypriority=false -blockmaxsize=16000000

# create more transactions without txfee
# optional
bitcoin-cli settxfee 0
```

### Usage:

1. edit conf/app.conf to be available
2. configure dispatch type in conf/app.conf named dispatch item

	- s2mTx: create a large of transaction(recommanded in first running)
	- m2sTx: aggregate less conin(recommanded in heavy listunspent load)
	- n2mTx: create more natural transactions with n inputs and m outputs randomly
	- s2sTx: aim to create transations recursively when there are abundantly spendable output

3. run

	```
	# compile in Linux
	go build -o rpc
	./rpc
	
	# compile in other platform, run in Linux
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o rpc
	scp ./rpc user@ip:/path/to/this reporisity
	./rpc
	```

### Result:

Access to [http://114.215.30.211:3002/](http://114.215.30.211:3002/) relaxly!

Blocks Height List(meeting expected result):

   `1219335` `121933` `1219456` `1219457` `1219759`

### Todo:

- [x] Create n2m transaction randomly
- [ ] Support dispatch functions via specified order in app.conf
- [ ] Fix unreachable items in data randomly after adding item

