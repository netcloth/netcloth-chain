#!/usr/bin/env bash

accs_file=$1
from=$2
passwd=$3
amt=$4
interval=$5

for acc in `cat $accs_file`
do
	echo $acc
	echo $passwd | nchcli send --from $from --to $acc --amount $amt -y
	sleep $interval
done
