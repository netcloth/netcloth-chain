acc_basename=$1
acc_number=$2

acc_number_max=$((acc_number+1))

for((i=1;i<$acc_number_max;i++))
do
	acc="$acc_basename$i"
	./account_add.sh $acc
done
