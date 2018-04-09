#!/bin/bash

set -e

if [[ $1 != "" ]]
then
    cd $1
    ./bitcoin-cli listunspent | grep amount | wc -l
else
    bitcoin-cli listunspent | grep amount | wc -l
fi
