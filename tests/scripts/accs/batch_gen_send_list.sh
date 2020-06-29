#!/usr/bin/env bash

accs_file_list=($@)
#echo ${accs_file_list[@]}

file_cnt=${#accs_file_list[*]}

if [ $file_cnt -lt 2 ]; then
	echo "at list 2 accs files"
	exit -1
fi

cnt=0
guard_num=$((file_cnt-1))
ff=
tf=

for((i=0;i<$file_cnt;i++))
do
	ff=${accs_file_list[$i]}

	tf_index=$((i+1))
	if [ $cnt -eq $guard_num ]; then
		tf_index=0
	fi

	tf=${accs_file_list[$tf_index]}

	echo "$ff-->$tf"
	bash gen_send_list.sh $ff $tf send_lists

	cnt=$((cnt+1))
done
