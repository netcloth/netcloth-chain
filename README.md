# netcloth-chain
An efficient blockchain network.

## QuickStart

### Build
set env
```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

build
```bash
# get source code
git clone https://github.com/NetCloth/netcloth-chain.git


# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
nchd help
nchcli help

```

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
nchd add-genesis-account $(nchcli keys show jack -a) 1000000000000unch,1000000000000stake
nchd add-genesis-account $(nchcli keys show alice -a) 1000000000000unch,1000000000000stake

# create validator
nchd gentx \
  --amount 1000000stake \
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
nchcli send --from nch13kvswghgd4gr9daymt2uykj5459luur4mrh4ef --to nch1cutxf4fe5twyqcegk88v7ll5aqhmkzcd68wg54 --amount 1stake
```