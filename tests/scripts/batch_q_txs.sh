#!/usr/bin/env bash

for tx in `cat txs`
do
nchcli q tx $tx
if [ $? -eq 1 ]; then
echo $tx >>badTxs
fi
done
