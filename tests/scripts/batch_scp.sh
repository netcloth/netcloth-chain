#!/usr/bin/env bash

machine_list=$1
what=$2
where=$3

for machine in `cat $machine_list`
do
scp $what $machine:$where
done
