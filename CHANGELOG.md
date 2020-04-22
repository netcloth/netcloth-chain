# Changelog

## [Unreleased]
### nchd

* fix #27 "--gas=auto" not available
* fix #29 op blockhash
* fix #32 op timestamp
* support #30 revert reason provided by the contract
* update cli interaction with contract txs
* update inflation model and related query api
* add gas_price_threshold param for tx gas price limit

### nchcli v1.0.4

#### [Features]

* update response when query account not exist 

* add distr alias for ```nchcli```

* support "--gas auto" for ```nchcli``` to calculate gas automatically

* add OOG logs in tx-receipt when out of gas

### nchd v1.0.4

#### [Features]

* export / import IPAL/C-IPAL state

#### [Bug Fixes]

* fix panic when C-IPAL data is null in genesis files
