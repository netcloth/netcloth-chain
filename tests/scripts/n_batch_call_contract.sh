n=$1
contract_addr=$2
for ((i=0;i<$n;i++))
do
	echo batch:$((i+1))
	bash batch_call_contract.sh $2
sleep 60
done
