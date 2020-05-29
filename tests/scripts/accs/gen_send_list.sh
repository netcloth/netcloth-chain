ff=$1
tf=$2
send_lists_dir=$3

mkdir -p $send_lists_dir

f_list=()
t_list=()

f_list+=(`cat $ff`)
t_list+=(`cat $tf`)

#echo ${f_list[*]}
#echo ${t_list[*]}

#echo ${#f_list[*]}
#echo ${#t_list[*]}

gen_file=$send_lists_dir/${ff}_to_${tf}
rm $gen_file

num=0
for f in ${f_list[*]}
do
	t=${t_list[$num]}
	num=$((num+1))
	echo $f:$t >> $gen_file
done
