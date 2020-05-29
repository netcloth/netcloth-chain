to_accs_file=$1
from_acc_file=$2
passwd=$3
amt=$4

batch_cnt=0
send_cnt=0
one_batch=()

do_send() {
for facc in `cat $from_acc_file`
do
	tacc=${one_batch[$send_cnt]}
	send_cnt=$((send_cnt+1))
	echo $facc-$tacc:$amt
	echo $passwd | nchcli send --from $facc --to $tacc --amount $amt -y
done
}

for acc in `cat $to_accs_file`
do
	one_batch+=($acc)
	batch_cnt=$((batch_cnt+1))
	if [ $batch_cnt -eq 20 ]; then
		do_send
		echo ok
		batch_cnt=0
		send_cnt=0
		unset one_batch
		sleep 6
	fi
done




