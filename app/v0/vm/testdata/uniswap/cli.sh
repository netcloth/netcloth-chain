#!/usr/bin/env bash

echo 11111111 | \
nchcli vm create --code_file=./uniswap.bc \
--from=$(nchcli keys show -a alice) \
--gas=9531375 \
-b block -y

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a alice) \
--contract_addr=nch1crgfeph6hl9z9hw0eevym6p26y85wy4gxle3fq \
--method=initializeFactory \
--abi_file="./uniswap.abi" \
--args="nch1w8x5hxmdwjz7pw0kwukmgaavnpu5q3e0churjj" \
--gas=98669 -b block -y

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a alice) \
--contract_addr=nch1crgfeph6hl9z9hw0eevym6p26y85wy4gxle3fq \
--method=createExchange \
--abi_file=./uniswap.abi \
--args="nch1w8x5hxmdwjz7pw0kwukmgaavnpu5q3e0churjj" \
--gas=4036200 -b block -y