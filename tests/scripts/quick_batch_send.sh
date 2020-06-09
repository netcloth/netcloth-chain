from_acc_file=$1
to_accs_file=$2
passwd=$3
amt=$4
interval=$5

one_batch_number=0

# get one_batch_number by from_acc_file
for facc in `cat $from_acc_file`
do
one_batch_number=$((one_batch_number+1))
done
echo $one_batch_number

one_batch=()

send_cnt=0
do_send() {
for facc in `cat $from_acc_file`
do
	tacc=${one_batch[$send_cnt]}
	send_cnt=$((send_cnt+1))
	echo "$facc-->$tacc:$amt"
	echo $passwd | nchcli send --from $facc --to $tacc --amount $amt -y &
done
}

batch_cnt=0
for acc in `cat $to_accs_file`
do
	one_batch+=($acc)
	batch_cnt=$((batch_cnt+1))
	if [ $batch_cnt -eq $one_batch_number ]; then
		do_send
		echo ok
		batch_cnt=0
		send_cnt=0
		unset one_batch
		sleep $interval
	fi
done
