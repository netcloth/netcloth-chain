#!/usr/bin/env bash

pkill nchd
pkill nchcli

sleep 1

nchd unsafe-reset-all

ssh n2 bash reset.sh
ssh n3 bash reset.sh
ssh n4 bash reset.sh
