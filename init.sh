#!/bin/bash
set -x

passwd="yourpassword"

rm -rf ~/.nchd
rm -rf ~/.nchcli

# set moniker and chain-id
nchd init mymoniker --chain-id nch-chain

# set up config for CLI
nchcli config chain-id nch-chain
nchcli config output json
nchcli config indent true
nchcli config trust-node true

# add keys
echo -e "${passwd}\n${passwd}\n" | nchcli keys add alice
echo -e "${passwd}\n${passwd}\n" | nchcli keys add bob
echo -e "${passwd}\n${passwd}\n" | nchcli keys add jack

# add genesis account
vesting_start_time=`date +%s`
vesting_end_time=`date -v+1d +%s`
nchd add-genesis-account $(nchcli keys show alice -a) 30000000000000000000pnch
nchd add-genesis-account $(nchcli keys show bob -a) 40000000000000000000pnch --vesting-amount 5000000000000000000pnch --vesting-start-time ${vesting_start_time} --vesting-end-time ${vesting_end_time}
nchd add-genesis-account $(nchcli keys show jack -a) 40000000000000000000pnch --vesting-amount 5000000000000000000pnch  --vesting-end-time ${vesting_end_time}

echo "${passwd}" | nchd gentx \
  --amount 1000000000000pnch \
  --commission-rate "0.10" \
  --commission-max-rate "0.20" \
  --commission-max-change-rate "0.10" \
  --pubkey $(nchd tendermint show-validator) \
  --name alice

# collect genesis tx
nchd collect-gentxs

# validate genesis file
nchd validate-genesis

# start the node
nchd start
