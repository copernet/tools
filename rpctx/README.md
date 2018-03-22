## Mine Big Size Blocks(BCH) Testing

This repository is under experiment. The current problem is that kafka is not suitable for transfering a big size message(up to 90 megabyte). Up to now, I have mined several blocks with about 16 megabyte size.

### Env

- Ubuntu 16.04.3 LTS
- Bitcoin-ABC v0.16.2

### Modified Source Code

```
# consensus.h
static const uint64_t DEFAULT_MAX_BLOCK_SIZE = 32 * ONE_MEGABYTE;
static const int64_t MAX_BLOCK_SIGOPS_PER_MB = 2000000;
static const int COINBASE_MATURITY = 1;  // optional

# policy.h
static const uint64_t DEFAULT_MAX_GENERATED_BLOCK_SIZE = 32 * ONE_MEGABYTE;
static const uint64_t DEFAULT_BLOCK_PRIORITY_PERCENTAGE = 0;
static const Amount DEFAULT_BLOCK_MIN_TX_FEE(0);
static const Amount DEFAULT_INCREMENTAL_RELAY_FEE(0);
static const Amount DUST_RELAY_TX_FEE(0);
```

> startup command: ./bitcoind -testnet -relaypriority=false

> bitcoin-cli settxfee 0    // after finishing startup

### Result

Access to [http://114.215.41.16:3002/](http://114.215.41.16:3002/) relaxly!

Blocks Height List(meeting expected result):

   `1219335` `121933` `1219456` `1219457` `1219759`

### Todo

- Create n2m transaction randomly
- Support dispatch functions via specified order in app.conf
- Fix unreachable items in data randomly after adding item

