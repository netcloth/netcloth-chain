contract_addr=$1
for acc in `cat accs1`
do
echo $acc
echo testtest | nchcli vm call --contract_addr=$contract_addr --abi_file ./payment/pay.abi --method=doTransfer --args 'nch169gj8454d02ld25wpsncca4xzqg2e02aqdau42' --amount 1000pnch --gas 30000000 -y --from $acc &
#sleep 0.03
done
