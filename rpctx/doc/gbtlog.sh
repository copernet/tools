#!/bin/bash
cd ~/.bitcoin-abc/testnet3
tail -f debug.log | grep -E 'CreateNewBlock|UpdateTip'

