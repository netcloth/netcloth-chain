send_file=$1
amt=$2
passwd=$3
cur_batch=$4
total_batch=$5
one_batch_cnt=`wc -l $send_file | awk '{print $1}'`

cnt=0
for send_pair in `cat $send_file`
do
	cnt=$((cnt+1))
	f=`echo $send_pair | awk -F ':' '{print $1}'`
	t=`echo $send_pair | awk -F ':' '{print $2}'`
	echo "cur_batch/total_batch:$cur_batch/$total_batch#$cnt/$one_batch_cnt"
	echo $passwd | nchcli send --from $f --to $t --amount $amt -y &
done
