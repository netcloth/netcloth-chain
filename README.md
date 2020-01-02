# netcloth-chain
An efficient blockchain network.

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
nchd add-genesis-account $(nchcli keys show jack -a) 1000000000000000000pnch
nchd add-genesis-account $(nchcli keys show alice -a) 1000000000000000000pnch

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
nchcli send --from $(nchcli keys show jack -a)  --to $(nchcli keys show alice -a) --amount 1pnch
```

query account
```
nchcli query account  $(nchcli keys show jack -a)
nchcli query account  $(nchcli keys show alice -a)
```

