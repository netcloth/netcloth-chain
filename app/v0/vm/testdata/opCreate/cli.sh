echo 11111111 | \
nchcli vm create --code_file=./a.bc \
--from=$(nchcli keys show -a alice) \
--gas=1000000 \
-b block -y

nchcli query vm call $(nchcli keys show -a alice) nch1crgfeph6hl9z9hw0eevym6p26y85wy4gxle3fq \
a "./a.abi"

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a alice) \
--contract_addr=nch1crgfeph6hl9z9hw0eevym6p26y85wy4gxle3fq \
--method=newAccount \
--abi_file="./a.abi" \
--gas=1000000 -b block -y
