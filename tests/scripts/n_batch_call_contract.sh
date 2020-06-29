#!/usr/bin/env bash

n=$1
contract_addr=$2
accs=$3
passwd=$4
interval=$5

for ((i=0;i<$n;i++))
do
	bash batch_call_contract.sh $contract_addr $accs $passwd $((i+1)) $interval
done
