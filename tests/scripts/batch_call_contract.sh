#!/usr/bin/env bash

contract_addr=$1
accs=$2
passwd=$3
batch=$4
interval=$5

batch_index=0

for acc in `cat $accs`
do

batch_index=$((batch_index+1))
echo $passwd | nchcli vm call --contract_addr=$contract_addr --abi_file ./payment/pay.abi --method=doTransfer --args 'nch169gj8454d02ld25wpsncca4xzqg2e02aqdau42' --amount 1000pnch --gas 30000000 -y --from $acc  &

echo "progress:$batch-$batch_index"
sleep $interval

done

