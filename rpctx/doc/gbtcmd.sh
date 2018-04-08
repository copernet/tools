#!/bin/bash

getBlockTemplate()
{
    while :
    do
         bitcoin-cli getblocktemplate > /dev/null
         sleep 5s
     done
}

getBlockTemplateWithParam()
{
    cd $1

    while :
    do
         ./bitcoin-cli getblocktemplate > /dev/null
         sleep 5s
     done
}

if [$1 -eq ""]
then
    getBlockTemplate
else
    getBlockTemplateWithParam
fi
