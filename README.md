<h1 align="center">NetCloth Chain</h1>

[![version](https://img.shields.io/github/tag/netcloth/netcloth-chain.svg)](https://github.com/netcloth/netcloth-chain/releases/latest)
[![license](https://img.shields.io/github/license/netcloth/netcloth-chain.svg)](https://github.com/netcloth/netcloth-chain/blob/master/LICENSE)
[![LoC](https://tokei.rs/b1/github/netcloth/netcloth-chain)](https://github.com/netcloth/netcloth-chain)
[![Go Report Card](https://goreportcard.com/badge/github.com/netcloth/netcloth-chain)](https://goreportcard.com/report/github.com/netcloth/netcloth-chain)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/fe23e8d8b57a47a782709085737aec4b)](https://app.codacy.com/gh/netcloth/netcloth-chain?utm_source=github.com&utm_medium=referral&utm_content=netcloth/netcloth-chain&utm_campaign=Badge_Grade_Dashboard)
[![CodeFactor](https://www.codefactor.io/repository/github/netcloth/netcloth-chain/badge)](https://www.codefactor.io/repository/github/netcloth/netcloth-chain)
[![codecov](https://codecov.io/gh/netcloth/netcloth-chain/branch/develop/graph/badge.svg)](https://codecov.io/gh/netcloth/netcloth-chain)
[![Build Status](https://travis-ci.com/netcloth/netcloth-chain.svg?branch=develop)](https://travis-ci.com/netcloth/netcloth-chain)

Welcome to the official Go implementation of the [NetCloth](https://www.netcloth.org) blockchain!

Founded in February 2019, NetCloth is a next-generation high-performance public chain network with the consensus of BPoS. NetCloth chain assists developers and users to create their own "personal network" by providing virtual machines (compatible with EVM), IPAL on-chain addressing protocol, native support for meta transactions, loop transactions and other infrastructures, and promotes the "personal network" to become The basic network unit and application service provider of the Web3.0. NetCloth assists users to truly master their personal data and establish a digital economic system based on data assets.

## QuickStart

### Install
Install nch from [here](https://github.com/netcloth/netcloth-chain/tree/develop/docs/install.md)

### Run
init
```
# Initialize configuration files and genesis file
nchd init local-nch --chain-id nch-chain

# Copy the `Address` output here and save it for later use 
nchcli keys add jack

# Copy the `Address` output here and save it for later use
nchcli keys add alice

# Add both accounts, with coins to the genesis file
nchd add-genesis-account $(nchcli keys show jack -a) 50000000000000000000pnch
nchd add-genesis-account $(nchcli keys show alice -a) 50000000000000000000pnch

# create validator
nchd gentx \
  --amount 1000000000000pnch \
  --commission-rate "0.10" \
  --commission-max-rate "0.20" \
  --commission-max-change-rate "0.10" \
  --pubkey $(nchd tendermint show-validator) \
  --name alice

# collect gentx
nchd collect-gentxs


# Configure your CLI to eliminate need for chain-id flag
nchcli config chain-id nch-chain
nchcli config output json
nchcli config indent true
nchcli config trust-node true
```

run nchd

```cassandraql
nchd start --log_level "*:debug" --trace
```

transfer asset
```cassandraql
# transfer asset
nchcli send --from $(nchcli keys show jack -a)  --to $(nchcli keys show alice -a) --amount 1000000000000pnch
```

query account
```
nchcli query account  $(nchcli keys show jack -a)
nchcli query account  $(nchcli keys show alice -a)
```

