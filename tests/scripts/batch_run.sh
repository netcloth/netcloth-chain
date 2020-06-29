#!/usr/bin/env bash

machine_list=$1
what="$2"
set -x

for machine in `cat $machine_list`
do
ssh $machine "$what"
done
